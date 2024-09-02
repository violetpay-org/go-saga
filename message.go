package saga

import (
	"encoding/json"
	"time"
)

type MessageConstructor func(Session) Message

// Message is a value object.
type Message interface {
	ID() string
	SessionID() string
	Trigger() string
	CreatedAt() time.Time

	// MarshalJSON returns the JSON encoding of the message.
	//
	// It must be implemented by the message struct that embeds the AbstractMessage.
	MarshalJSON() ([]byte, error)
}

// AbstractMessage is a value object that represents a message.
// It contains the common fields of a message.
// If you want to create a new message, you should embed this struct.
// and must implement the MarshalJSON method.
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

func (m AbstractMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string    `json:"id"`
		SessionID string    `json:"sessionID"`
		Trigger   string    `json:"trigger"`
		CreatedAt time.Time `json:"createdAt"`
	}{
		ID:        m.id,
		SessionID: m.sessionID,
		Trigger:   m.trigger,
		CreatedAt: m.createdAt,
	})
}

type AbstractMessageRepository[Tx TxContext] interface {
	AbstractMessageLoadRepository
	SaveMessage(message Message) Executable[Tx]
	SaveMessages(messages []Message) Executable[Tx]
	SaveDeadLetter(message Message) Executable[Tx]
	SaveDeadLetters(message []Message) Executable[Tx]
	DeleteMessage(message Message) Executable[Tx]
	DeleteMessages(messages []Message) Executable[Tx]
	DeleteDeadLetter(message Message) Executable[Tx]
	DeleteDeadLetters(messages []Message) Executable[Tx]
}

type AbstractMessageLoadRepository interface {
	GetMessagesFromOutbox(batchSize int) ([]Message, error)
	GetMessagesFromDeadLetter(batchSize int) ([]Message, error)
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
