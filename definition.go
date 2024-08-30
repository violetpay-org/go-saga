package saga

type Definition struct {
	steps []Step
}

func newDefinition(steps []Step) Definition {
	return Definition{
		steps: steps,
	}
}
