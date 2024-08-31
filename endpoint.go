package saga

type Endpoint[Tx TxContext] struct {
	commandChannel     ChannelName
	commandConstructor MessageConstructor
	commandRepository  AbstractMessageRepository[Tx]

	successResChannel          ChannelName
	successResponseConstructor MessageConstructor
	failureResChannel          ChannelName
	failureResponseConstructor MessageConstructor
}

func NewEndpoint[Tx TxContext](
	commandChannel ChannelName,
	commandConstructor MessageConstructor,
	commandRepository AbstractMessageRepository[Tx],
	successResChannel ChannelName,
	successResponseConstructor MessageConstructor,
	failureResChannel ChannelName,
	failureResponseConstructor MessageConstructor,
) Endpoint[Tx] {
	return Endpoint[Tx]{
		commandChannel:             commandChannel,
		commandConstructor:         commandConstructor,
		commandRepository:          commandRepository,
		successResChannel:          successResChannel,
		successResponseConstructor: successResponseConstructor,
		failureResChannel:          failureResChannel,
		failureResponseConstructor: failureResponseConstructor,
	}
}

func (e Endpoint[Tx]) CommandChannel() ChannelName {
	return e.commandChannel
}

func (e Endpoint[Tx]) CommandConstructor() MessageConstructor {
	return e.commandConstructor
}

func (e Endpoint[Tx]) CommandRepository() AbstractMessageRepository[Tx] {
	return e.commandRepository
}

func (e Endpoint[Tx]) SuccessResChannel() ChannelName {
	return e.successResChannel
}

func (e Endpoint[Tx]) SuccessResponseConstructor() MessageConstructor {
	return e.successResponseConstructor
}

func (e Endpoint[Tx]) FailureResChannel() ChannelName {
	return e.failureResChannel
}

func (e Endpoint[Tx]) FailureResponseConstructor() MessageConstructor {
	return e.failureResponseConstructor
}

type ExecutablePreparer[Tx TxContext] func(Session) (Executable[Tx], error)

type LocalEndpoint[Tx TxContext] struct {
	successResChannel          ChannelName
	successResponseConstructor MessageConstructor
	successResRepository       AbstractMessageRepository[Tx]

	failureResChannel          ChannelName
	failureResponseConstructor MessageConstructor
	failureResRepository       AbstractMessageRepository[Tx]

	handler ExecutablePreparer[Tx]
}

func NewLocalEndpoint[Tx TxContext](
	successResChannel ChannelName,
	successResponseConstructor MessageConstructor,
	successResRepository AbstractMessageRepository[Tx],
	failureResChannel ChannelName,
	failureResponseConstructor MessageConstructor,
	failureResRepository AbstractMessageRepository[Tx],
	handler ExecutablePreparer[Tx],
) LocalEndpoint[Tx] {
	return LocalEndpoint[Tx]{
		successResChannel:          successResChannel,
		successResponseConstructor: successResponseConstructor,
		successResRepository:       successResRepository,
		failureResChannel:          failureResChannel,
		failureResponseConstructor: failureResponseConstructor,
		failureResRepository:       failureResRepository,
		handler:                    handler,
	}
}

func (e LocalEndpoint[Tx]) SuccessResChannel() ChannelName {
	return e.successResChannel
}

func (e LocalEndpoint[Tx]) SuccessResponseConstructor() MessageConstructor {
	return e.successResponseConstructor
}

func (e LocalEndpoint[Tx]) SuccessResRepository() AbstractMessageRepository[Tx] {
	return e.successResRepository
}

func (e LocalEndpoint[Tx]) FailureResChannel() ChannelName {
	return e.failureResChannel
}

func (e LocalEndpoint[Tx]) FailureResponseConstructor() MessageConstructor {
	return e.failureResponseConstructor
}

func (e LocalEndpoint[Tx]) FailureResRepository() AbstractMessageRepository[Tx] {
	return e.failureResRepository
}

func (e LocalEndpoint[Tx]) handle(session Session) (cmd Executable[Tx], err error) {
	cmd, err = e.handler(session)
	return
}
