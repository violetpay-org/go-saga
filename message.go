package saga

import (
	"time"
)

type MessageConstructor[S Session, M Message] func(S) M

func convertMessage[S Session, M Message](constructor MessageConstructor[S, M]) MessageConstructor[Session, Message] {
	return func(s Session) Message {
		return constructor(s.(S))
	}
}

// Message is a value object.
type Message interface {
	ID() string
	SessionID() string
	Trigger() string
	CreatedAt() time.Time
}

func NewAbstractMessage(id, sessionID, trigger string) AbstractMessage {
	return AbstractMessage{
		id:        id,
		sessionID: sessionID,
		trigger:   trigger,
		createdAt: time.Now(),
	}
}

func NewAbstractMessageWithTime(id, sessionID, trigger string, createdAt time.Time) AbstractMessage {
	return AbstractMessage{
		id:        id,
		sessionID: sessionID,
		trigger:   trigger,
		createdAt: createdAt,
	}
}

// AbstractMessage is a value object that represents a message.
// It contains the common fields of a message.
// If you want to create a new message, you should embed this struct.
type AbstractMessage struct {
	id        string
	sessionID string
	trigger   string
	createdAt time.Time
}

func (m AbstractMessage) ID() string {
	return m.id
}

func (m AbstractMessage) SessionID() string {
	return m.sessionID
}

func (m AbstractMessage) Trigger() string {
	return m.trigger
}

func (m AbstractMessage) CreatedAt() time.Time {
	return m.createdAt
}

func ConvertMessageRepository[M Message, Tx TxContext](repository AbstractMessageRepository[M, Tx]) AbstractMessageRepository[Message, Tx] {
	getMessagesFromOutbox := func(batchSize int) ([]Message, error) {
		ms, err := repository.GetMessagesFromOutbox(batchSize)
		if err != nil {
			return nil, err
		}
		var messages []Message
		for _, m := range ms {
			messages = append(messages, m)
		}
		return messages, nil
	}

	getMessagesFromDeadLetter := func(batchSize int) ([]Message, error) {
		ms, err := repository.GetMessagesFromDeadLetter(batchSize)
		if err != nil {
			return nil, err
		}
		var messages []Message
		for _, m := range ms {
			messages = append(messages, m)
		}
		return messages, nil
	}

	saveMessages := func(messages []Message) Executable[Tx] {
		var ms []M
		for _, m := range messages {
			ms = append(ms, m.(M))
		}
		return repository.SaveMessages(ms)
	}

	saveDeadLetters := func(messages []Message) Executable[Tx] {
		var ms []M
		for _, m := range messages {
			ms = append(ms, m.(M))
		}
		return repository.SaveDeadLetters(ms)
	}

	deleteMessages := func(messages []Message) Executable[Tx] {
		var ms []M
		for _, m := range messages {
			ms = append(ms, m.(M))
		}
		return repository.DeleteMessages(ms)
	}

	deleteDeadLetters := func(messages []Message) Executable[Tx] {
		var ms []M
		for _, m := range messages {
			ms = append(ms, m.(M))
		}
		return repository.DeleteDeadLetters(ms)
	}

	return messageRepository[Tx]{
		saveMessage:               func(m Message) Executable[Tx] { return repository.SaveMessage(m.(M)) },
		saveMessages:              saveMessages,
		saveDeadLetter:            func(m Message) Executable[Tx] { return repository.SaveDeadLetter(m.(M)) },
		saveDeadLetters:           saveDeadLetters,
		deleteMessage:             func(m Message) Executable[Tx] { return repository.DeleteMessage(m.(M)) },
		deleteMessages:            deleteMessages,
		deleteDeadLetter:          func(m Message) Executable[Tx] { return repository.DeleteDeadLetter(m.(M)) },
		deleteDeadLetters:         deleteDeadLetters,
		getMessagesFromOutbox:     getMessagesFromOutbox,
		getMessagesFromDeadLetter: getMessagesFromDeadLetter,
	}
}

type AbstractMessageRepository[M Message, Tx TxContext] interface {
	AbstractMessageLoadRepository[M]
	SaveMessage(message M) Executable[Tx]
	SaveMessages(messages []M) Executable[Tx]
	SaveDeadLetter(message M) Executable[Tx]
	SaveDeadLetters(message []M) Executable[Tx]
	DeleteMessage(message M) Executable[Tx]
	DeleteMessages(messages []M) Executable[Tx]
	DeleteDeadLetter(message M) Executable[Tx]
	DeleteDeadLetters(messages []M) Executable[Tx]
}

type AbstractMessageLoadRepository[M Message] interface {
	GetMessagesFromOutbox(batchSize int) ([]M, error)
	GetMessagesFromDeadLetter(batchSize int) ([]M, error)
}

type messageRepository[Tx TxContext] struct {
	saveMessage               func(Message) Executable[Tx]
	saveMessages              func([]Message) Executable[Tx]
	saveDeadLetter            func(Message) Executable[Tx]
	saveDeadLetters           func([]Message) Executable[Tx]
	deleteMessage             func(Message) Executable[Tx]
	deleteMessages            func([]Message) Executable[Tx]
	deleteDeadLetter          func(Message) Executable[Tx]
	deleteDeadLetters         func([]Message) Executable[Tx]
	getMessagesFromOutbox     func(int) ([]Message, error)
	getMessagesFromDeadLetter func(int) ([]Message, error)
}

func (r messageRepository[Tx]) SaveMessage(message Message) Executable[Tx] {
	return r.saveMessage(message)
}

func (r messageRepository[Tx]) SaveMessages(messages []Message) Executable[Tx] {
	return r.saveMessages(messages)
}

func (r messageRepository[Tx]) SaveDeadLetter(message Message) Executable[Tx] {
	return r.saveDeadLetter(message)
}

func (r messageRepository[Tx]) SaveDeadLetters(messages []Message) Executable[Tx] {
	return r.saveDeadLetters(messages)
}

func (r messageRepository[Tx]) DeleteMessage(message Message) Executable[Tx] {
	return r.deleteMessage(message)
}

func (r messageRepository[Tx]) DeleteMessages(messages []Message) Executable[Tx] {
	return r.deleteMessages(messages)
}

func (r messageRepository[Tx]) DeleteDeadLetter(message Message) Executable[Tx] {
	return r.deleteDeadLetter(message)
}

func (r messageRepository[Tx]) DeleteDeadLetters(messages []Message) Executable[Tx] {
	return r.deleteDeadLetters(messages)
}

func (r messageRepository[Tx]) GetMessagesFromOutbox(batchSize int) ([]Message, error) {
	return r.getMessagesFromOutbox(batchSize)
}

func (r messageRepository[Tx]) GetMessagesFromDeadLetter(batchSize int) ([]Message, error) {
	return r.getMessagesFromDeadLetter(batchSize)
}

// messagePacket is a value object that represents a AbstractMessage packet.
// The packet contains the Message and the origin channel of the AbstractMessage.
type messagePacket struct {
	origin ChannelName
	Message
}

func (m messagePacket) Origin() ChannelName {
	return m.origin
}

func (m messagePacket) Payload() Message {
	return m.Message
}

func newMessagePacket(origin ChannelName, message Message) messagePacket {
	return messagePacket{
		Message: message,
		origin:  origin,
	}
}

type Command struct {
}

type Response struct {
}
