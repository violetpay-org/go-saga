package saga

type Definition struct {
	steps []Step
}

func newDefinition(steps []Step) Definition {
	return Definition{
		steps: steps,
	}
}

func (d Definition) FirstStep() Step {
	if len(d.steps) == 0 {
		return nil
	}

	return d.steps[0]
}

func (d Definition) FindStep(name string) Step {
	for _, step := range d.steps {
		if step.Name() == name {
			return step
		}
	}

	return nil
}

func (d Definition) Exists(step Step) bool {
	for _, s := range d.steps {
		if s.Name() == step.Name() {
			return true
		}
	}

	return false
}

func (d Definition) NextStep(step Step) Step {
	for i, s := range d.steps {
		if s.Name() == step.Name() {
			if i+1 < len(d.steps) {
				return d.steps[i+1]
			}
		}
	}

	return nil
}

func (d Definition) PrevStep(step Step) Step {
	for i, s := range d.steps {
		if s.Name() == step.Name() {
			if i-1 >= 0 {
				return d.steps[i-1]
			}
		}
	}

	return nil
}
