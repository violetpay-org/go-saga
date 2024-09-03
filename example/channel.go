package main

import (
	"errors"
	"github.com/violetpay-org/go-saga"
	"github.com/violetpay-org/go-saga/messageRelayer"
)

const (
	ExampleSuccessChannelName = "ExampleSuccessChannel"
	ExampleFailureChannelName = "ExampleFailureChannel"
	ExampleCommandChannelName = "ExampleCommandChannel"
)

var ExampleSuccessChannel = saga.NewChannel[ExampleMessage, ExampleTxContext](ExampleSuccessChannelName, registry, exampleSuccessResponseRepository)
var ExampleFailureChannel = saga.NewChannel[ExampleMessage, ExampleTxContext](ExampleFailureChannelName, registry, exampleFailureResponseRepository) // repo ?
var ExampleCommandChannel = messageRelayer.NewChannel[ExampleMessage, ExampleTxContext](
	ExampleCommandChannelName,
	registry,
	exampleCommandRepository,
	func(message saga.Message) error {
		return ExampleSuccessChannel.Send(message)
	},
)
var AlwaysFailCommandChannel = messageRelayer.NewChannel[ExampleMessage, ExampleTxContext](
	"AlwaysFailCommandChannel",
	registry,
	exampleCommandRepository,
	func(message saga.Message) error {
		return errors.New("AlwaysFailCommandChannel failed")
	},
)

var channelRegistry = messageRelayer.NewChannelRegistry[ExampleTxContext]()

func init() {
	var err error
	err = channelRegistry.Register(ExampleSuccessChannel)
	err = channelRegistry.Register(ExampleFailureChannel)
	err = channelRegistry.Register(ExampleCommandChannel)
	err = channelRegistry.Register(AlwaysFailCommandChannel)

	if err != nil {
		panic(err)
	}
}
