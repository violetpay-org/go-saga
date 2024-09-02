package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/violetpay-org/go-saga"
	"sync"
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
	return &ExampleMessageRepository{}
}

type ExampleMessageRepository struct {
	outbox     sync.Map
	deadLetter sync.Map
}

func (e *ExampleMessageRepository) GetMessagesFromOutbox(batchSize int) ([]saga.Message, error) {
	var outbox []saga.Message

	e.outbox.Range(func(key, value interface{}) bool {
		outbox = append(outbox, value.(saga.Message))
		return true
	})

	if batchSize > len(outbox) {
		return outbox, nil
	}

	return outbox[:batchSize], nil
}

func (e *ExampleMessageRepository) GetMessagesFromDeadLetter(batchSize int) ([]saga.Message, error) {
	var deadLetter []saga.Message

	e.outbox.Range(func(key, value interface{}) bool {
		deadLetter = append(deadLetter, value.(saga.Message))
		return true
	})

	if batchSize > len(deadLetter) {
		return deadLetter, nil
	}

	return deadLetter[:batchSize], nil
}

func (e *ExampleMessageRepository) SaveMessage(message saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.outbox.Store(message.ID(), message)
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
		e.deadLetter.Store(message.ID(), message)
		return nil
	}
}

func (e *ExampleMessageRepository) SaveDeadLetters(messages []saga.Message) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.SaveDeadLetter(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) DeleteMessage(message saga.Message) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.outbox.Delete(message.ID())
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
		e.deadLetter.Delete(message.ID())
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
