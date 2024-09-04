package saga

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepBuilder(t *testing.T) {
	repo := newMockAbstractMessageRepository()
	messageConstructor := func(session Session) mockMessage {
		return newMockMessage()
	}

	endpoint := NewEndpoint[Session, mockMessage, mockMessage, mockMessage, mockTxContext](
		"",
		messageConstructor,
		repo,
		"",
		messageConstructor,
		"",
		messageConstructor,
	)

	handler := func(session Session) (Executable[mockTxContext], error) {
		return nil, nil
	}

	localEndpoint := NewLocalEndpoint[Session, mockMessage, mockMessage, mockTxContext](
		"",
		messageConstructor,
		repo,
		"",
		messageConstructor,
		repo,
		handler,
	)

	var builder StepBuilder[mockTxContext]
	builder = StepBuilder[mockTxContext]{
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
			Invoke(endpoint).
			WithCompensation(endpoint).
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
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
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
			Invoke(endpoint).
			WithCompensation(endpoint).
			Retry().
			Step(step2Name).
			Invoke(endpoint).
			WithCompensation(endpoint).
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
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
			Retry().
			Step(step2Name).
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
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
			Invoke(endpoint).
			WithCompensation(endpoint).
			Retry().
			Step(step2Name).
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
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
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
			Retry().
			Step(step2Name).
			Invoke(endpoint).
			WithCompensation(endpoint).
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
			Invoke(endpoint).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(localEndpoint).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(endpoint).
			Retry().
			Step(step2Name).
			Invoke(endpoint).
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
			Invoke(endpoint).
			WithCompensation(endpoint).
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(endpoint).
			WithCompensation(endpoint).
			Step(step2Name).
			Invoke(endpoint).
			WithCompensation(endpoint).
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
			Invoke(endpoint).
			WithCompensation(endpoint).
			Retry().
			Step(step2Name).
			LocalInvoke(localEndpoint).
			Step(step3Name).
			Invoke(endpoint).
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
			LocalInvoke(localEndpoint).
			WithLocalCompensation(localEndpoint).
			Retry().
			Step(step2Name).
			Invoke(endpoint).
			Step(step3Name).
			LocalInvoke(localEndpoint).
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

}
