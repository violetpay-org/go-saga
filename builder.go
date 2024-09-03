package saga

type stepBuilder[Tx TxContext] interface {
	Step(name string) invocableBuild[Tx]
	Build() Definition
}

type invocableBuild[Tx TxContext] interface {
	Invoke(endpoint Endpoint[Tx]) invocationOptionBuild[Tx]
	LocalInvoke(endpoint LocalEndpoint[Tx]) localInvocationOptionBuild[Tx]
}

type invocationOptionBuild[Tx TxContext] interface {
	stepBuilder[Tx]
	compensableBuild[Tx]
	retryableBuild[Tx]
}

type compensableBuild[Tx TxContext] interface {
	WithCompensation(endpoint Endpoint[Tx]) invokeOptionWithoutCompensationBuild[Tx]
}

type localInvocationOptionBuild[Tx TxContext] interface {
	stepBuilder[Tx]
	localCompensableBuild[Tx]
	retryableBuild[Tx]
}

type localCompensableBuild[Tx TxContext] interface {
	WithLocalCompensation(endpoint LocalEndpoint[Tx]) invokeOptionWithoutCompensationBuild[Tx]
}

type invokeOptionWithoutCompensationBuild[Tx TxContext] interface {
	stepBuilder[Tx]
	retryableBuild[Tx]
}

type retryableBuild[Tx TxContext] interface {
	Retry() stepBuilder[Tx]
}

type StepBuilder[Tx TxContext] struct {
	steps           []Step
	currentStep     Step
	currentStepName string
}

func NewStepBuilder[Tx TxContext]() stepBuilder[Tx] {
	return &StepBuilder[Tx]{
		steps: make([]Step, 0),
	}
}

func (b *StepBuilder[Tx]) Invoke(endpoint Endpoint[Tx]) invocationOptionBuild[Tx] {
	b.currentStep = newRemoteStep(b.currentStepName, endpoint)

	return b
}

func (b *StepBuilder[Tx]) LocalInvoke(endpoint LocalEndpoint[Tx]) localInvocationOptionBuild[Tx] {
	b.currentStep = newLocalStep(b.currentStepName, endpoint)

	return b
}

func (b *StepBuilder[Tx]) Build() Definition {
	if b.currentStep != nil {
		b.steps = append(b.steps, b.currentStep)
	}
	defer b.cleanUp()

	return newDefinition(b.steps)
}

func (b *StepBuilder[Tx]) cleanUp() {
	b.steps = make([]Step, 0)
	b.currentStep = nil
	b.currentStepName = ""
}

func (b *StepBuilder[Tx]) WithCompensation(endpoint Endpoint[Tx]) invokeOptionWithoutCompensationBuild[Tx] {
	b.currentStep = newRemoteStepWithCompensation(b.currentStep.(remoteStep[Tx]), endpoint)
	return b
}

func (b *StepBuilder[Tx]) WithLocalCompensation(endpoint LocalEndpoint[Tx]) invokeOptionWithoutCompensationBuild[Tx] {
	b.currentStep = newLocalStepWithCompensation(b.currentStep.(localStep[Tx]), endpoint)
	return b
}

func (b *StepBuilder[Tx]) Retry() stepBuilder[Tx] {
	var newRet Step
	switch b.currentStep.(type) {
	case remoteStep[Tx]:
		newRet = newRemoteStepWithRetry(b.currentStep.(remoteStep[Tx]))
	case localStep[Tx]:
		newRet = newLocalStepWithRetry(b.currentStep.(localStep[Tx]))
	default:
		panic("Unknown step type")
	}

	b.currentStep = newRet
	return b
}

func (b *StepBuilder[Tx]) Step(name string) invocableBuild[Tx] {
	if b.currentStep != nil {
		b.steps = append(b.steps, b.currentStep)
	}

	b.currentStep = nil
	b.currentStepName = name

	return b
}
