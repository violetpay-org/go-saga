package main

import (
	"github.com/google/uuid"
	"github.com/violetpay-org/go-saga"
	"sync"
)

var exampleCommandRepository = NewExampleMessageRepository()
var exampleSuccessResponseRepository = NewExampleMessageRepository()
var exampleFailureResponseRepository = NewExampleMessageRepository()

func ExampleMessageConstructor(session *ExampleSession) ExampleMessage {
	return ExampleMessage{
		AbstractMessage: saga.NewAbstractMessage(
			uuid.New().String(),
			session.ID(),
			"Triggered by test",
		),
		exampleField: session.exampleField,
	}
}

type ExampleMessage struct {
	saga.AbstractMessage
	exampleField string
}

func NewExampleMessageRepository() *ExampleMessageRepository {
	return &ExampleMessageRepository{}
}

type ExampleMessageRepository struct {
	outbox     sync.Map
	deadLetter sync.Map
}

func (e *ExampleMessageRepository) GetMessagesFromOutbox(batchSize int) ([]ExampleMessage, error) {
	var outbox []ExampleMessage

	e.outbox.Range(func(key, value interface{}) bool {
		outbox = append(outbox, value.(ExampleMessage))
		return true
	})

	if batchSize > len(outbox) {
		return outbox, nil
	}

	return outbox[:batchSize], nil
}

func (e *ExampleMessageRepository) GetMessagesFromDeadLetter(batchSize int) ([]ExampleMessage, error) {
	var deadLetter []ExampleMessage

	e.deadLetter.Range(func(key, value interface{}) bool {
		deadLetter = append(deadLetter, value.(ExampleMessage))
		return true
	})

	if batchSize > len(deadLetter) {
		return deadLetter, nil
	}

	return deadLetter[:batchSize], nil
}

func (e *ExampleMessageRepository) SaveMessage(message ExampleMessage) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.outbox.Store(message.ID(), message)
		return nil
	}
}

func (e *ExampleMessageRepository) SaveMessages(messages []ExampleMessage) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.SaveMessage(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) SaveDeadLetter(message ExampleMessage) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.deadLetter.Store(message.ID(), message)
		return nil
	}
}

func (e *ExampleMessageRepository) SaveDeadLetters(messages []ExampleMessage) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.SaveDeadLetter(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) DeleteMessage(message ExampleMessage) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.outbox.Delete(message.ID())
		return nil
	}
}

func (e *ExampleMessageRepository) DeleteMessages(messages []ExampleMessage) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.DeleteMessage(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) DeleteDeadLetter(message ExampleMessage) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.deadLetter.Delete(message.ID())
		return nil
	}
}

func (e *ExampleMessageRepository) DeleteDeadLetters(messages []ExampleMessage) saga.Executable[ExampleTxContext] {
	executables := make([]saga.Executable[ExampleTxContext], 0)
	for _, msg := range messages {
		executables = append(executables, e.DeleteDeadLetter(msg))
	}

	return saga.CombineExecutables(executables...)
}

func (e *ExampleMessageRepository) clear() {
	e.outbox = sync.Map{}
	e.deadLetter = sync.Map{}
}
