package saga

type stepBuilder interface {
	Step(name string) invocableBuild
	Build() Definition
}

type invocableBuild interface {
	Invoke(endpoint Endpoint) invocationOptionBuild
	LocalInvoke(endpoint LocalEndpoint) localInvocationOptionBuild
}

type invocationOptionBuild interface {
	stepBuilder
	compensableBuild
	retryableBuild
}

type compensableBuild interface {
	WithCompensation(endpoint Endpoint) invokeOptionWithoutCompensationBuild
}

type localInvocationOptionBuild interface {
	stepBuilder
	localCompensableBuild
	retryableBuild
}

type localCompensableBuild interface {
	WithLocalCompensation(endpoint LocalEndpoint) invokeOptionWithoutCompensationBuild
}

type invokeOptionWithoutCompensationBuild interface {
	stepBuilder
	retryableBuild
}

type retryableBuild interface {
	Retry() stepBuilder
}

type StepBuilder struct {
	steps           []Step
	currentStep     Step
	currentStepName string
}

func NewStepBuilder() stepBuilder {
	return &StepBuilder{
		steps: make([]Step, 0),
	}
}

func (b *StepBuilder) Invoke(endpoint Endpoint) invocationOptionBuild {
	b.currentStep = newRemoteStep(b.currentStepName, endpoint)

	return b
}

func (b *StepBuilder) LocalInvoke(endpoint LocalEndpoint) localInvocationOptionBuild {
	b.currentStep = newLocalStep(b.currentStepName, endpoint)

	return b
}

func (b *StepBuilder) Build() Definition {
	if b.currentStep != nil {
		b.steps = append(b.steps, b.currentStep)
	}
	defer b.cleanUp()

	if len(b.steps) == 0 {
		panic("No steps defined, but trying to build a definition")
	}

	return newDefinition(b.steps)
}

func (b *StepBuilder) cleanUp() {
	b.steps = make([]Step, 0)
	b.currentStep = nil
	b.currentStepName = ""
}

func (b *StepBuilder) WithCompensation(endpoint Endpoint) invokeOptionWithoutCompensationBuild {
	b.currentStep = newRemoteStepWithCompensation(b.currentStep.(remoteStep), endpoint)
	return b
}

func (b *StepBuilder) WithLocalCompensation(endpoint LocalEndpoint) invokeOptionWithoutCompensationBuild {
	b.currentStep = newLocalStepWithCompensation(b.currentStep.(localStep), endpoint)
	return b
}

func (b *StepBuilder) Retry() stepBuilder {
	var newRet Step
	switch b.currentStep.(type) {
	case remoteStep:
		newRet = newRemoteStepWithRetry(b.currentStep.(remoteStep))
	case localStep:
		newRet = newLocalStepWithRetry(b.currentStep.(localStep))
	}

	b.currentStep = newRet
	return b
}

func (b *StepBuilder) Step(name string) invocableBuild {
	if b.currentStep != nil {
		b.steps = append(b.steps, b.currentStep)
	}

	b.currentStep = nil
	b.currentStepName = name

	return b
}
