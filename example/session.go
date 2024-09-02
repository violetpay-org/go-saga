package main

import (
	"errors"
	"fmt"
	"github.com/violetpay-org/go-saga"
	"sync"
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
	return &ExampleSessionRepository{}
}

type ExampleSessionRepository struct {
	sessions sync.Map
}

func (e *ExampleSessionRepository) Load(id string) (*ExampleSession, error) {
	sess, ok := e.sessions.Load(id)
	if !ok {
		return nil, errors.New("session not found")
	}

	session := sess.(ExampleSession)

	return &session, nil
}

func (e *ExampleSessionRepository) Save(sess *ExampleSession) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		fmt.Println("Saving session", sess.ID(), sess.State(), sess.IsPending(), sess.currentStep.Name())
		e.sessions.Store(sess.ID(), *sess)
		return nil
	}
}

func (e *ExampleSessionRepository) Delete(sess *ExampleSession) saga.Executable[ExampleTxContext] {
	return func(ctx ExampleTxContext) error {
		e.sessions.Delete(sess.ID())
		return nil
	}
}
