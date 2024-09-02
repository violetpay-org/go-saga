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
	builder := saga.NewStepBuilder[ExampleTxContext]()
	def := builder.Step("ExampleStep1").Invoke(ExampleEndpoint).Step("ExampleStep2").Invoke(ExampleEndpoint).Build()

	e.Saga = saga.NewSaga[*ExampleSession, ExampleTxContext](
		"ExampleSaga",
		def,
		exampleSessionFactory,
		exampleSessionRepository,
	)
}

func (e *ExampleSaga) ApplySchemaTo(registry *saga.Registry[ExampleTxContext]) {
	e.buildSaga()
	saga.RegisterSagaTo(registry, e.Saga)
}
