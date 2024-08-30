package saga

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepBuilder(t *testing.T) {
	var builder StepBuilder
	builder = StepBuilder{
		steps: make([]Step, 0),
	}

	step1Name := "step1"
	step2Name := "step2"
	step3Name := "step3"
	//var endpoint Endpoint
	//var localEndpoint LocalEndpoint

	t.Run("Build", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			Invoke(nil).
			WithCompensation(nil).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())
	})

	t.Run("Build with local endpoint", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())
	})

	t.Run("Build with multiple steps", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			Invoke(nil).
			WithCompensation(nil).
			Retry().
			Step(step2Name).
			Invoke(nil).
			WithCompensation(nil).
			Retry().
			Build()

		assert.Equal(t, 2, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.True(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.True(t, def.steps[1].MustBeCompleted())
	})

	t.Run("Build with multiple steps with local endpoint", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Retry().
			Step(step2Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Retry().
			Build()

		assert.Equal(t, 2, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.True(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.True(t, def.steps[1].MustBeCompleted())
	})

	t.Run("Build with multiple steps with mixed endpoint", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			Invoke(nil).
			WithCompensation(nil).
			Retry().
			Step(step2Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Retry().
			Build()

		assert.Equal(t, 2, len(def.steps))
		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.True(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.True(t, def.steps[1].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Retry().
			Step(step2Name).
			Invoke(nil).
			WithCompensation(nil).
			Retry().
			Build()

		assert.Equal(t, 2, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.True(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.True(t, def.steps[1].MustBeCompleted())
	})

	t.Run("Build with no options", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			Invoke(nil).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(nil).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(nil).
			Retry().
			Step(step2Name).
			Invoke(nil).
			Retry().
			Build()

		assert.Equal(t, 2, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.False(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.True(t, def.steps[1].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(nil).
			WithCompensation(nil).
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(nil).
			WithCompensation(nil).
			Step(step2Name).
			Invoke(nil).
			WithCompensation(nil).
			Build()

		assert.Equal(t, 2, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.True(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.False(t, def.steps[1].MustBeCompleted())
	})

	t.Run("Build 3 steps with mixed endpoint", func(t *testing.T) {
		def := builder.
			Step(step1Name).
			Invoke(nil).
			WithCompensation(nil).
			Retry().
			Step(step2Name).
			LocalInvoke(nil).
			Step(step3Name).
			Invoke(nil).
			Retry().
			Build()

		assert.Equal(t, 3, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.False(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.False(t, def.steps[1].MustBeCompleted())

		assert.Equal(t, step3Name, def.steps[2].Name())
		assert.False(t, def.steps[2].IsCompensable())
		assert.True(t, def.steps[2].IsInvocable())
		assert.True(t, def.steps[2].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(nil).
			WithLocalCompensation(nil).
			Retry().
			Step(step2Name).
			Invoke(nil).
			Step(step3Name).
			LocalInvoke(nil).
			Retry().
			Build()

		assert.Equal(t, 3, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		assert.Equal(t, step2Name, def.steps[1].Name())
		assert.False(t, def.steps[1].IsCompensable())
		assert.True(t, def.steps[1].IsInvocable())
		assert.False(t, def.steps[1].MustBeCompleted())

		assert.Equal(t, step3Name, def.steps[2].Name())
		assert.False(t, def.steps[2].IsCompensable())
		assert.True(t, def.steps[2].IsInvocable())
		assert.True(t, def.steps[2].MustBeCompleted())
	})

	t.Run("Build with no steps", func(t *testing.T) {
		assert.Panics(t, func() {
			builder.Build()
		})
	})
}
