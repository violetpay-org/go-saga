package saga

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testRepo struct {
}

func (r *testRepo) GetMessagesFromOutbox(batchSize int) ([]Message, error) {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) GetMessagesFromDeadLetter(batchSize int) ([]Message, error) {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) SaveMessage(message Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) SaveMessages(messages []Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) SaveDeadLetter(message Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) SaveDeadLetters(message []Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) DeleteMessage(message Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) DeleteMessages(messages []Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) DeleteDeadLetter(message Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

func (r *testRepo) DeleteDeadLetters(message Message) Executable[tx] {
	//TODO implement me
	panic("implement me")
}

type tx interface{}

var (
	msgConstructor = func(session Session) Message {
		return nil
	}

	txEndpoint = NewEndpoint[tx](
		"",
		msgConstructor,
		&testRepo{},
		"",
		msgConstructor,
		"",
		msgConstructor,
	)

	txLocalEndpoint = NewLocalEndpoint[tx](
		"",
		msgConstructor,
		&testRepo{},
		"",
		msgConstructor,
		&testRepo{},
		func(session Session) (Executable[tx], error) {
			return nil, nil
		},
	)
)

func TestStepBuilder(t *testing.T) {
	var builder StepBuilder[tx]

	builder = StepBuilder[tx]{
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
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
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
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
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
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
			Retry().
			Step(step2Name).
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
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
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
			Retry().
			Step(step2Name).
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
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
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
			Retry().
			Step(step2Name).
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
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
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
			Retry().
			Step(step2Name).
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
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
			Invoke(txEndpoint).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(txLocalEndpoint).
			Retry().
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.False(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.True(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(txEndpoint).
			Retry().
			Step(step2Name).
			Invoke(txEndpoint).
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
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
			Build()

		assert.Equal(t, 1, len(def.steps))

		assert.Equal(t, step1Name, def.steps[0].Name())
		assert.True(t, def.steps[0].IsCompensable())
		assert.True(t, def.steps[0].IsInvocable())
		assert.False(t, def.steps[0].MustBeCompleted())

		def = builder.
			Step(step1Name).
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
			Step(step2Name).
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
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
			Invoke(txEndpoint).
			WithCompensation(txEndpoint).
			Retry().
			Step(step2Name).
			LocalInvoke(txLocalEndpoint).
			Step(step3Name).
			Invoke(txEndpoint).
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
			LocalInvoke(txLocalEndpoint).
			WithLocalCompensation(txLocalEndpoint).
			Retry().
			Step(step2Name).
			Invoke(txEndpoint).
			Step(step3Name).
			LocalInvoke(txLocalEndpoint).
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
