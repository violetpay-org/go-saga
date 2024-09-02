package saga

type Step interface {
	// Name returns the name of the step.
	Name() string

	// IsCompensable returns true if the step is compensable.
	IsCompensable() bool

	// IsInvocable returns true if the step is invocable.
	IsInvocable() bool

	// MustBeCompleted returns true if the invocation action must be completed. So If it is true, something will be retried until it is completed.
	MustBeCompleted() bool
}

func newRemoteStep[Tx TxContext](name string, endpoint Endpoint[Tx]) remoteStep[Tx] {
	return remoteStep[Tx]{
		name:           name,
		invocation:     newRemoteInvocationAction(endpoint),
		invokeEndpoint: endpoint,
		retry:          false,
	}
}

func newRemoteStepWithCompensation[Tx TxContext](step remoteStep[Tx], endpoint Endpoint[Tx]) remoteStep[Tx] {
	return remoteStep[Tx]{
		name:           step.name,
		invocation:     step.invocation,
		invokeEndpoint: step.invokeEndpoint,
		compensation:   newRemoteCompensationAction(endpoint),
		compEndpoint:   endpoint,
		retry:          step.retry,
	}
}

func newRemoteStepWithRetry[Tx TxContext](step remoteStep[Tx]) remoteStep[Tx] {
	return remoteStep[Tx]{
		name:           step.name,
		invocation:     step.invocation,
		invokeEndpoint: step.invokeEndpoint,
		compensation:   step.compensation,
		compEndpoint:   step.compEndpoint,
		retry:          true,
	}
}

type remoteStep[Tx TxContext] struct {
	name string

	invocation     invokeAction[Tx]
	invokeEndpoint Endpoint[Tx]
	compensation   compensateAction[Tx]
	compEndpoint   Endpoint[Tx]

	retry bool
}

func (s remoteStep[Tx]) Name() string {
	return s.name
}

func (s remoteStep[Tx]) IsCompensable() bool {
	return s.compensation != nil
}

func (s remoteStep[Tx]) IsInvocable() bool {
	return s.invocation != nil
}

func (s remoteStep[Tx]) MustBeCompleted() bool {
	return s.retry == true
}

func (s remoteStep[Tx]) SetRetry(retry bool) remoteStep[Tx] {
	s.retry = retry
	return s
}

func newRemoteCompensationAction[Tx TxContext](endpoint Endpoint[Tx]) compensateAction[Tx] {
	return func(s Session) Executable[Tx] {
		command := endpoint.CommandConstructor()(s)
		return endpoint.CommandRepository().SaveMessage(command)
	}
}

type compensateAction[Tx TxContext] func(Session) Executable[Tx]

func newRemoteInvocationAction[Tx TxContext](endpoint Endpoint[Tx]) invokeAction[Tx] {
	return func(s Session) Executable[Tx] {
		command := endpoint.CommandConstructor()(s)
		return endpoint.CommandRepository().SaveMessage(command)
	}
}

type invokeAction[Tx TxContext] func(Session) Executable[Tx]

func newLocalStep[Tx TxContext](name string, endpoint LocalEndpoint[Tx]) localStep[Tx] {
	return localStep[Tx]{
		name: name,

		invocation:     newLocalInvokeAction(endpoint),
		invokeEndpoint: endpoint,
		retry:          false,
	}
}

func newLocalStepWithCompensation[Tx TxContext](step localStep[Tx], endpoint LocalEndpoint[Tx]) localStep[Tx] {
	return localStep[Tx]{
		name:           step.name,
		invocation:     step.invocation,
		invokeEndpoint: step.invokeEndpoint,
		compensation:   newLocalCompensateAction(endpoint),
		compEndpoint:   endpoint,
		retry:          step.retry,
	}
}

func newLocalStepWithRetry[Tx TxContext](step localStep[Tx]) localStep[Tx] {
	return localStep[Tx]{
		name:           step.name,
		invocation:     step.invocation,
		invokeEndpoint: step.invokeEndpoint,
		compensation:   step.compensation,
		compEndpoint:   step.compEndpoint,
		retry:          true,
	}
}

type localStep[Tx TxContext] struct {
	name string

	invocation     localInvokeAction[Tx]
	invokeEndpoint LocalEndpoint[Tx]
	compensation   localCompensateAction[Tx]
	compEndpoint   LocalEndpoint[Tx]

	retry bool
}

func (s localStep[Tx]) Name() string {
	return s.name
}

func (s localStep[Tx]) IsCompensable() bool {
	return s.compensation != nil
}

func (s localStep[Tx]) IsInvocable() bool {
	return s.invocation != nil
}

func (s localStep[Tx]) MustBeCompleted() bool {
	return s.retry == true
}

func (s localStep[Tx]) SetRetry(retry bool) localStep[Tx] {
	s.retry = retry
	return s
}

func newLocalCompensateAction[Tx TxContext](endpoint LocalEndpoint[Tx]) localCompensateAction[Tx] {
	action := func(s Session) (Executable[Tx], error) {
		cmd, err := endpoint.handle(s)
		if err != nil {
			msg := endpoint.FailureResponseConstructor()(s)
			return endpoint.FailureResRepository().SaveMessage(msg), nil
		}

		msg := endpoint.SuccessResponseConstructor()(s)
		cmd2 := endpoint.SuccessResRepository().SaveMessage(msg)
		return CombineExecutables(cmd, cmd2), nil
	}

	return action
}

type localCompensateAction[Tx TxContext] func(Session) (Executable[Tx], error)

func newLocalInvokeAction[Tx TxContext](endpoint LocalEndpoint[Tx]) localInvokeAction[Tx] {
	return func(s Session) (Executable[Tx], error) {
		cmd, err := endpoint.handle(s)
		if err != nil {
			msg := endpoint.FailureResponseConstructor()(s)
			return endpoint.FailureResRepository().SaveMessage(msg), nil
		}

		msg := endpoint.SuccessResponseConstructor()(s)
		cmd2 := endpoint.SuccessResRepository().SaveMessage(msg)
		return CombineExecutables(cmd, cmd2), nil
	}
}

type localInvokeAction[Tx TxContext] func(Session) (Executable[Tx], error)
