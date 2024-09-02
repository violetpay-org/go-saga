package main

import (
	"context"
	"github.com/violetpay-org/go-saga"
)

var UnitOfWorkFactory = saga.NewUnitOfWorkFactory[ExampleTxContext](ExampleTxHandler{})

type ExampleTxContext struct{}

type ExampleTxHandler struct{}

func (e ExampleTxHandler) BeginTx(ctx context.Context) (tx ExampleTxContext, error error) {
	return ExampleTxContext{}, nil
}

func (e ExampleTxHandler) Commit(ctx ExampleTxContext) error {
	return nil
}

func (e ExampleTxHandler) Rollback(ctx ExampleTxContext) error {
	return nil
}
