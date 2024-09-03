package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/violetpay-org/go-saga"
	"github.com/violetpay-org/go-saga/messageRelayer"
	"testing"
)

func TestOrchestrator(t *testing.T) {
	successRepository := exampleSuccessResponseRepository
	failureRepository := exampleFailureResponseRepository
	commandRepository := exampleCommandRepository
	sessionRepository := exampleSessionRepository
	var exampleSaga saga.Saga[*ExampleSession, ExampleTxContext]

	CleanUp := func(t *testing.T) {
		registry = saga.NewRegistry(orchestrator)
		successRepository.clear()
		failureRepository.clear()
		commandRepository.clear()
		sessionRepository.clear()

		ExampleSuccessChannel = saga.NewChannel[ExampleMessage, ExampleTxContext](ExampleSuccessChannelName, registry, exampleSuccessResponseRepository)
		ExampleFailureChannel = saga.NewChannel[ExampleMessage, ExampleTxContext](ExampleFailureChannelName, registry, exampleFailureResponseRepository) // repo ?
		ExampleCommandChannel = messageRelayer.NewChannel[ExampleMessage, ExampleTxContext](
			ExampleCommandChannelName,
			registry,
			exampleCommandRepository,
			func(message saga.Message) error {
				return ExampleSuccessChannel.Send(message)
			},
		)
		AlwaysFailCommandChannel = messageRelayer.NewChannel[ExampleMessage, ExampleTxContext](
			"AlwaysFailCommandChannel",
			registry,
			exampleCommandRepository,
			func(message saga.Message) error {
				return errors.New("AlwaysFailCommandChannel failed")
			},
		)

		channelRegistry = messageRelayer.NewChannelRegistry[ExampleTxContext]()
		err := channelRegistry.RegisterChannel(ExampleSuccessChannel)
		err = channelRegistry.RegisterChannel(ExampleFailureChannel)
		err = channelRegistry.RegisterChannel(ExampleCommandChannel)
		err = channelRegistry.RegisterChannel(AlwaysFailCommandChannel)
		assert.Nil(t, err)

		exampleSaga = saga.Saga[*ExampleSession, ExampleTxContext]{}
	}

	builder := saga.NewStepBuilder[ExampleTxContext]()

	buildSagaAndRegister := func(def saga.Definition) {
		exampleSaga = saga.NewSaga[*ExampleSession, ExampleTxContext](
			"ExampleSaga",
			def,
			exampleSessionFactory,
			exampleSessionRepository,
		)

		err := saga.RegisterSagaTo(registry, exampleSaga)
		if err != nil {
			panic(err)
		}
	}

	t.Run("should be able to build and register saga", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(builder.Build())

		assert.True(t, registry.HasSaga(exampleSaga.Name()))
	})

	t.Run("should be throw because of no first step", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(builder.Build())

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.NotNil(t, err)
	})

	t.Run("should got session after saga started", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				Invoke(ExampleEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
	})

	t.Run("should be return error when register duplicated saga", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				Invoke(ExampleEndpoint).
				Build(),
		)

		err := saga.RegisterSagaTo(registry, exampleSaga)
		assert.NotNil(t, err)
	})

	t.Run("should set session is pending when start saga", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				Invoke(ExampleEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateCommon, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())
	})

	t.Run("should be complete saga when all steps are done", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				LocalInvoke(ExampleLocalEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		err = messageRelayer.New(1, channelRegistry, UnitOfWorkFactory).Execute()
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.False(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateCompleted, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())
	})

	t.Run("should be fail saga because of local endpoint handle() called", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(builder.
			Step("ExampleStep1").
			LocalInvoke(ExampleAlwaysFailingLocalEndpoint).
			Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		err = messageRelayer.New(1, channelRegistry, UnitOfWorkFactory).Execute()
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.False(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateFailed, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())
	})

	t.Run("should be retry saga step when endpoint fails", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				LocalInvoke(ExampleAlwaysFailingLocalEndpoint).
				Retry().
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		relayer := messageRelayer.New(1, channelRegistry, UnitOfWorkFactory)
		// Consume first step
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateIsRetrying, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())

		// Consume second time, retry
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err = sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateIsRetrying, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())
	})

	t.Run("should be produced success response when a local endpoint is done", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				LocalInvoke(ExampleLocalEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		outbox, err := exampleSuccessResponseRepository.GetMessagesFromOutbox(10)
		assert.Nil(t, err)

		deadLetter, err := exampleSuccessResponseRepository.GetMessagesFromDeadLetter(10)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(outbox))
		assert.Equal(t, 0, len(deadLetter))

		outbox, err = exampleFailureResponseRepository.GetMessagesFromOutbox(10)
		assert.Nil(t, err)

		deadLetter, err = exampleFailureResponseRepository.GetMessagesFromDeadLetter(10)
		assert.Nil(t, err)

		assert.Equal(t, 0, len(outbox))
		assert.Equal(t, 0, len(deadLetter))
	})

	t.Run("should be produced failure response when a local endpoint is failed", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				LocalInvoke(ExampleAlwaysFailingLocalEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		outbox, err := exampleFailureResponseRepository.GetMessagesFromOutbox(10)
		assert.Nil(t, err)

		deadLetter, err := exampleFailureResponseRepository.GetMessagesFromDeadLetter(10)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(outbox))
		assert.Equal(t, 0, len(deadLetter))

		outbox, err = exampleSuccessResponseRepository.GetMessagesFromOutbox(10)
		assert.Nil(t, err)

		deadLetter, err = exampleSuccessResponseRepository.GetMessagesFromDeadLetter(10)

		assert.Equal(t, 0, len(outbox))
		assert.Equal(t, 0, len(deadLetter))
	})

	t.Run("should be return dead session error when consumed same message twice", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				Invoke(ExampleEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		all, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		err = ExampleSuccessChannel.Send(ExampleMessage{
			AbstractMessage: saga.NewAbstractMessage(
				uuid.New().String(),
				all[0].ID(),
				"Triggered by test",
			),
			exampleField: "",
		})
		assert.Nil(t, err)

		// Consume first time
		err = messageRelayer.New(1, channelRegistry, UnitOfWorkFactory).Execute()
		assert.Nil(t, err)

		err = ExampleSuccessChannel.Send(ExampleMessage{
			AbstractMessage: saga.NewAbstractMessage(
				uuid.New().String(),
				all[0].ID(),
				"Triggered by test",
			),
			exampleField: "",
		})
		assert.NotNil(t, err)
	})

	t.Run("should be return error when start saga with empty name", func(t *testing.T) {
		CleanUp(t)

		err := registry.StartSaga("", map[string]interface{}{})
		assert.NotNil(t, err)
	})

	t.Run("should be return error when start saga with empty name", func(t *testing.T) {
		CleanUp(t)

		err := registry.StartSaga("Test", nil)
		assert.NotNil(t, err)
	})

	t.Run("should be success multiple steps", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				LocalInvoke(ExampleLocalEndpoint).
				Step("ExampleStep2").
				LocalInvoke(ExampleLocalEndpoint).
				Step("ExampleStep3").
				LocalInvoke(ExampleLocalEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		relayer := messageRelayer.New(1, channelRegistry, UnitOfWorkFactory)

		// Consume first step
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateCommon, sessions[0].State())
		assert.Equal(t, "ExampleStep2", sessions[0].CurrentStep().Name())

		// Consume second step
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err = sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateCommon, sessions[0].State())
		assert.Equal(t, "ExampleStep3", sessions[0].CurrentStep().Name())

		// Consume third step
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err = sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.False(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateCompleted, sessions[0].State())
		assert.Equal(t, "ExampleStep3", sessions[0].CurrentStep().Name())
	})

	t.Run("should be compensate when a local endpoint is failed", func(t *testing.T) {
		CleanUp(t)

		buildSagaAndRegister(
			builder.
				Step("ExampleStep1").
				LocalInvoke(ExampleLocalEndpoint).
				WithLocalCompensation(ExampleLocalEndpoint).
				Step("ExampleStep2").
				LocalInvoke(ExampleAlwaysFailingLocalEndpoint).
				Build(),
		)

		err := registry.StartSaga(exampleSaga.Name(), map[string]interface{}{})
		assert.Nil(t, err)

		// Consume first step
		relayer := messageRelayer.New(1, channelRegistry, UnitOfWorkFactory)
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err := sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateCommon, sessions[0].State())
		assert.Equal(t, "ExampleStep2", sessions[0].CurrentStep().Name())

		// Consume second step but failed, back to prev step
		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err = sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.True(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateIsCompensating, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())

		// Consume first compensation step
		err = ExampleSuccessChannel.Send(ExampleMessage{
			AbstractMessage: saga.NewAbstractMessage(
				uuid.New().String(),
				sessions[0].ID(),
				"Triggered by test",
			),
		})
		assert.Nil(t, err)

		err = relayer.Execute()
		assert.Nil(t, err)

		sessions, err = sessionRepository.loadAll()
		assert.Nil(t, err)

		assert.Equal(t, 1, len(sessions))
		assert.False(t, sessions[0].IsPending())
		assert.Equal(t, saga.StateFailed, sessions[0].State())
		assert.Equal(t, "ExampleStep1", sessions[0].CurrentStep().Name())
	})
}
