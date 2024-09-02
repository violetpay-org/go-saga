package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/violetpay-org/go-saga"
	"time"
)

var exampleCommandRepository = NewExampleMessageRepository()
var exampleSuccessResponseRepository = NewExampleMessageRepository()
var exampleFailureResponseRepository = NewExampleMessageRepository()

func ExampleMessageConstructor(session saga.Session) saga.Message {
	return &ExampleMessage{
		AbstractMessage: saga.NewAbstractMessage(
			uuid.New().String(),
			session.ID(),
			"Triggered by test",
		),
		exampleField: session.(*ExampleSession).exampleField,
	}
}

type ExampleMessage struct {
	saga.AbstractMessage
	exampleField string
}

func (m *ExampleMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           string    `json:"id"`
		SessionID    string    `json:"sessionID"`
		Trigger      string    `json:"trigger"`
		CreatedAt    time.Time `json:"createdAt"`
		ExampleField string    `json:"exampleField"`
	}{
		ID:           m.ID(),
		SessionID:    m.SessionID(),
		Trigger:      m.Trigger(),
		CreatedAt:    m.CreatedAt(),
		ExampleField: m.exampleField,
	})
}

func NewExampleMessageRepository() *ExampleMessageRepository {
	return &ExampleMessageRepository{
		outbox:     make(map[string]saga.Message),
		deadLetter: make(map[string]saga.Message),
	}
}

type ExampleMessageRepository struct {
	outbox     map[string]saga.Message
	deadLetter map[string]saga.Message
}

func (e *ExampleMessageRepository) GetMessagesFromOutbox(batchSize int) ([]saga.Message, error) {
	var outbox []saga.Message
	for _, message := range e.outbox {
		outbox = append(outbox, message)
	}

	if batchSize > len(outbox) {
		return outbox, nil
	}

	return outbox[:batchSize], nil
}

func (e *ExampleMessageRepository) GetMessagesFromDeadLetter(batchSize int) ([]saga.Message, error) {
	var deadLetter []saga.Message
	for _, message := range e.deadLetter {
		deadLetter = append(deadLetter, message)
	}

	if batchSize > len(deadLetter) {
		return deadLetter, nil
	}

	return deadLetter[:batchSize], nil
}

func (e *ExampleMessageRepository) SaveMessage(message saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.outbox[message.ID()] = message
		return nil
	}
}

func (e *ExampleMessageRepository) SaveMessages(messages []saga.Message) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.SaveMessage(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) SaveDeadLetter(message saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.deadLetter[message.ID()] = message
		return nil
	}
}

func (e *ExampleMessageRepository) SaveDeadLetters(message []saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		for _, msg := range message {
			e.deadLetter[msg.ID()] = msg
		}
		return nil
	}
}

func (e *ExampleMessageRepository) DeleteMessage(message saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		delete(e.outbox, message.ID())
		return nil
	}
}

func (e *ExampleMessageRepository) DeleteMessages(messages []saga.Message) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.DeleteMessage(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) DeleteDeadLetter(message saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		delete(e.deadLetter, message.ID())
		return nil
	}
}

func (e *ExampleMessageRepository) DeleteDeadLetters(messages []saga.Message) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.DeleteDeadLetter(msg))
	}

	return saga.CombineExecutables(executables...)
}
