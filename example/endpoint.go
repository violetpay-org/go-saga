package main

import (
	"errors"
	"github.com/violetpay-org/go-saga"
)

var ExampleEndpoint = saga.NewEndpoint[
	*ExampleSession,
	ExampleMessage, ExampleMessage, ExampleMessage,
	ExampleTxContext,
](
	ExampleCommandChannelName,
	ExampleMessageConstructor,
	exampleCommandRepository,
	ExampleSuccessChannelName,
	ExampleMessageConstructor,
	ExampleFailureChannelName,
	ExampleMessageConstructor,
)

var ExampleLocalEndpoint = saga.NewLocalEndpoint[
	*ExampleSession,
	ExampleMessage, ExampleMessage,
	ExampleTxContext,
](
	ExampleSuccessChannelName,
	ExampleMessageConstructor,
	exampleSuccessResponseRepository,
	ExampleFailureChannelName,
	ExampleMessageConstructor,
	exampleFailureResponseRepository,
	func(session saga.Session) (saga.Executable[ExampleTxContext], error) {
		return func(ctx ExampleTxContext) error {
			return nil
		}, nil
	},
)

var ExampleAlwaysFailingLocalEndpoint = saga.NewLocalEndpoint[
	*ExampleSession,
	ExampleMessage, ExampleMessage,
	ExampleTxContext,
](
	ExampleSuccessChannelName,
	ExampleMessageConstructor,
	exampleSuccessResponseRepository,
	ExampleFailureChannelName,
	ExampleMessageConstructor,
	exampleFailureResponseRepository,
	func(session saga.Session) (saga.Executable[ExampleTxContext], error) {
		return func(ctx ExampleTxContext) error {
			return nil
		}, errors.New("failed because always failing endpoint called")
	},
)
