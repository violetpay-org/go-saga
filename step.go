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

func newRemoteStep(name string, endpoint Endpoint) remoteStep {
	return remoteStep{
		name: name,

		invocation: newRemoteInvocationAction(endpoint),
		retry:      false,
	}
}

func newRemoteStepWithCompensation(step remoteStep, endpoint Endpoint) remoteStep {
	return remoteStep{
		name:         step.name,
		invocation:   step.invocation,
		compensation: newRemoteCompensationAction(endpoint),
		retry:        step.retry,
	}
}

func newRemoteStepWithRetry(step remoteStep) remoteStep {
	return remoteStep{
		name:         step.name,
		invocation:   step.invocation,
		compensation: step.compensation,
		retry:        true,
	}
}

type remoteStep struct {
	name string

	invocation   *invokeAction
	compensation *compensateAction
	retry        bool
}

func (s remoteStep) Name() string {
	return s.name
}

func (s remoteStep) IsCompensable() bool {
	return s.compensation != nil
}

func (s remoteStep) IsInvocable() bool {
	return s.invocation != nil
}

func (s remoteStep) MustBeCompleted() bool {
	return s.retry == true
}

func (s remoteStep) SetRetry(retry bool) remoteStep {
	s.retry = retry
	return s
}

func newRemoteCompensationAction(endpoint Endpoint) *compensateAction {
	return &compensateAction{}
}

type compensateAction struct {
}

func newRemoteInvocationAction(endpoint Endpoint) *invokeAction {
	return &invokeAction{}
}

type invokeAction struct {
}

func newLocalStep(name string, endpoint LocalEndpoint) localStep {
	return localStep{
		name: name,

		invocation: newLocalInvokeAction(endpoint),
		retry:      false,
	}
}

func newLocalStepWithCompensation(step localStep, endpoint LocalEndpoint) localStep {
	return localStep{
		name:         step.name,
		invocation:   step.invocation,
		compensation: newLocalCompensateAction(endpoint),
		retry:        step.retry,
	}
}

func newLocalStepWithRetry(step localStep) localStep {
	return localStep{
		name:         step.name,
		invocation:   step.invocation,
		compensation: step.compensation,
		retry:        true,
	}
}

type localStep struct {
	name string

	compensation *localCompensateAction
	invocation   *localInvokeAction
	retry        bool
}

func (s localStep) Name() string {
	return s.name
}

func (s localStep) IsCompensable() bool {
	return s.compensation != nil
}

func (s localStep) IsInvocable() bool {
	return s.invocation != nil
}

func (s localStep) MustBeCompleted() bool {
	return s.retry == true
}

func (s localStep) SetRetry(retry bool) localStep {
	s.retry = retry
	return s
}

func newLocalCompensateAction(endpoint Endpoint) *localCompensateAction {
	return &localCompensateAction{}
}

type localCompensateAction struct {
}

func newLocalInvokeAction(endpoint Endpoint) *localInvokeAction {
	return &localInvokeAction{}
}

type localInvokeAction struct {
}
