package main

import "github.com/violetpay-org/go-saga"

var ExampleEndpoint = saga.NewEndpoint[*ExampleSession, ExampleMessage, ExampleTxContext](
	ExampleCommandChannelName,
	ExampleMessageConstructor,
	exampleCommandRepository,
	ExampleSuccessChannelName,
	ExampleMessageConstructor,
	ExampleFailureChannelName,
	ExampleMessageConstructor,
)
