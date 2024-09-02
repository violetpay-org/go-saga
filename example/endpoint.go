package main

import "github.com/violetpay-org/go-saga"

var ExampleEndpoint = saga.NewEndpoint[ExampleTxContext](
	ExampleCommandChannelName,
	ExampleMessageConstructor,
	exampleCommandRepository,
	ExampleSuccessChannelName,
	ExampleMessageConstructor,
	ExampleFailureChannelName,
	ExampleMessageConstructor,
)
