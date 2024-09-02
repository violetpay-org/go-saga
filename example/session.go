package main

import (
	"errors"
	"github.com/violetpay-org/go-saga"
	"log"
)

var exampleSessionRepository = NewExampleSessionRepository()
var exampleSessionFactory = func(data map[string]interface{}) *ExampleSession {
	return &ExampleSession{
		id: data["id"].(string),
	}
}

type ExampleSession struct {
	id           string
	currentStep  saga.Step
	pending      bool
	state        saga.State
	exampleField string
}

func (e *ExampleSession) ID() string {
	return e.id
}

func (e *ExampleSession) CurrentStep() saga.Step {
	return e.currentStep

}

func (e *ExampleSession) UpdateCurrentStep(step saga.Step) error {
	e.currentStep = step
	return nil
}

func (e *ExampleSession) IsPending() bool {
	return e.pending
}

func (e *ExampleSession) SetPending(pending bool) {
	e.pending = pending
}

func (e *ExampleSession) State() saga.State {
	return e.state
}

func (e *ExampleSession) SetState(state saga.State) {
	e.state = state
}

func NewExampleSessionRepository() *ExampleSessionRepository {
	return &ExampleSessionRepository{
		sessions: make(map[string]ExampleSession),
	}
}

type ExampleSessionRepository struct {
	sessions map[string]ExampleSession
}

func (e *ExampleSessionRepository) Load(id string) (*ExampleSession, error) {
	sess, ok := e.sessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}
	return &sess, nil
}

func (e *ExampleSessionRepository) Save(sess *ExampleSession) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.sessions[sess.ID()] = *sess
		return nil
	}
}

func (e *ExampleSessionRepository) Delete(sess *ExampleSession) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		log.Print("Deleting session", sess.ID())
		delete(e.sessions, sess.ID())
		return nil
	}
}
