package saga

type Endpoint[Tx TxContext] struct {
	commandChannel     ChannelName
	commandConstructor MessageConstructor[Session, Message]
	commandRepository  AbstractMessageRepository[Tx]

	successResChannel          ChannelName
	successResponseConstructor MessageConstructor[Session, Message]
	failureResChannel          ChannelName
	failureResponseConstructor MessageConstructor[Session, Message]
}

func NewEndpoint[S Session, C Message, SRes Message, FRes Message, Tx TxContext](
	commandChannel ChannelName,
	commandConstructor MessageConstructor[S, C],
	commandRepository AbstractMessageRepository[Tx],
	successResChannel ChannelName,
	successResponseConstructor MessageConstructor[S, SRes],
	failureResChannel ChannelName,
	failureResponseConstructor MessageConstructor[S, FRes],
) Endpoint[Tx] {
	return Endpoint[Tx]{
		commandChannel:             commandChannel,
		commandConstructor:         convertMessage(commandConstructor),
		commandRepository:          commandRepository,
		successResChannel:          successResChannel,
		successResponseConstructor: convertMessage(successResponseConstructor),
		failureResChannel:          failureResChannel,
		failureResponseConstructor: convertMessage(failureResponseConstructor),
	}
}

func (e Endpoint[Tx]) CommandChannel() ChannelName {
	return e.commandChannel
}

func (e Endpoint[Tx]) CommandConstructor() MessageConstructor[Session, Message] {
	return e.commandConstructor
}

func (e Endpoint[Tx]) CommandRepository() AbstractMessageRepository[Tx] {
	return e.commandRepository
}

func (e Endpoint[Tx]) SuccessResChannel() ChannelName {
	return e.successResChannel
}

func (e Endpoint[Tx]) SuccessResponseConstructor() MessageConstructor[Session, Message] {
	return e.successResponseConstructor
}

func (e Endpoint[Tx]) FailureResChannel() ChannelName {
	return e.failureResChannel
}

func (e Endpoint[Tx]) FailureResponseConstructor() MessageConstructor[Session, Message] {
	return e.failureResponseConstructor
}

type ExecutablePreparer[Tx TxContext] func(Session) (Executable[Tx], error)

type LocalEndpoint[Tx TxContext] struct {
	successResChannel          ChannelName
	successResponseConstructor MessageConstructor[Session, Message]
	successResRepository       AbstractMessageRepository[Tx]

	failureResChannel          ChannelName
	failureResponseConstructor MessageConstructor[Session, Message]
	failureResRepository       AbstractMessageRepository[Tx]

	handler ExecutablePreparer[Tx]
}

func NewLocalEndpoint[S Session, SRes Message, FRes Message, Tx TxContext](
	successResChannel ChannelName,
	successResponseConstructor MessageConstructor[S, SRes],
	successResRepository AbstractMessageRepository[Tx],
	failureResChannel ChannelName,
	failureResponseConstructor MessageConstructor[S, FRes],
	failureResRepository AbstractMessageRepository[Tx],
	handler ExecutablePreparer[Tx],
) LocalEndpoint[Tx] {
	return LocalEndpoint[Tx]{
		successResChannel:          successResChannel,
		successResponseConstructor: convertMessage(successResponseConstructor),
		successResRepository:       successResRepository,
		failureResChannel:          failureResChannel,
		failureResponseConstructor: convertMessage(failureResponseConstructor),
		failureResRepository:       failureResRepository,
		handler:                    handler,
	}
}

func (e LocalEndpoint[Tx]) SuccessResChannel() ChannelName {
	return e.successResChannel
}

func (e LocalEndpoint[Tx]) SuccessResponseConstructor() MessageConstructor[Session, Message] {
	return e.successResponseConstructor
}

func (e LocalEndpoint[Tx]) SuccessResRepository() AbstractMessageRepository[Tx] {
	return e.successResRepository
}

func (e LocalEndpoint[Tx]) FailureResChannel() ChannelName {
	return e.failureResChannel
}

func (e LocalEndpoint[Tx]) FailureResponseConstructor() MessageConstructor[Session, Message] {
	return e.failureResponseConstructor
}

func (e LocalEndpoint[Tx]) FailureResRepository() AbstractMessageRepository[Tx] {
	return e.failureResRepository
}

func (e LocalEndpoint[Tx]) handle(session Session) (cmd Executable[Tx], err error) {
	cmd, err = e.handler(session)
	return
}
