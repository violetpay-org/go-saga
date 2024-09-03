package main

import (
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

var channelRegistry = messageRelayer.NewChannelRegistry[ExampleTxContext]()

func init() {
	var err error
	err = channelRegistry.RegisterChannel(ExampleSuccessChannel)
	err = channelRegistry.RegisterChannel(ExampleFailureChannel)
	err = channelRegistry.RegisterChannel(ExampleCommandChannel)

	if err != nil {
		panic(err)
	}
}
