package main

import (
	"github.com/violetpay-org/go-saga"
)

func NewExampleSaga() *ExampleSaga {
	return &ExampleSaga{}
}

type ExampleSaga struct {
	saga.Saga[*ExampleSession, ExampleTxContext]
}

func (e *ExampleSaga) buildSaga() {
	def := saga.NewStepBuilder[ExampleTxContext]().
		Step("ExampleStep1").
		LocalInvoke(ExampleLocalEndpoint).
		Step("ExampleStep2").
		LocalInvoke(ExampleLocalEndpoint).
		Build()

	e.Saga = saga.NewSaga[*ExampleSession, ExampleTxContext](
		"ExampleSaga",
		def,
		exampleSessionFactory,
		exampleSessionRepository,
	)
}

func (e *ExampleSaga) ApplySchemaTo(registry *saga.Registry[ExampleTxContext]) error {
	e.buildSaga()
	return saga.RegisterSagaTo(registry, e.Saga)
}
