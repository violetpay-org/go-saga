package main

import "github.com/violetpay-org/go-saga"

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
