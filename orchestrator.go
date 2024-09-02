package saga

import (
	"context"
	"errors"
	"log"
)

type Orchestrator[Tx TxContext] interface {
	Orchestrate(saga Saga[Session, Tx], packet messagePacket) error
	StartSaga(saga Saga[Session, Tx], sessionArgs map[string]interface{}) error
}

func NewOrchestrator[Tx TxContext](uowFactory UnitOfWorkFactory[Tx]) Orchestrator[Tx] {
	return &orchestrator[Tx]{
		uowFactory: uowFactory,
	}
}

type orchestrator[Tx TxContext] struct {
	uowFactory UnitOfWorkFactory[Tx]
}

func (o *orchestrator[Tx]) StartSaga(saga Saga[Session, Tx], sessionArgs map[string]interface{}) error {
	var uow *UnitOfWork[Tx]
	var err error

	sagaSession := saga.createSession(sessionArgs)
	if sagaSession == nil {
		return errors.New("session is nil")
	}

	if sagaSession.ID() == "" {
		return errors.New("session ID is empty")
	}

	sagaDef := saga.Definition()

	firstStep := sagaDef.FirstStep()
	if firstStep == nil {
		return errors.New("saga has no first step")
	}

	log.Print("Updating current step")

	err = sagaSession.UpdateCurrentStep(firstStep)
	if err != nil {
		return err
	}

	uow, err = o.uowFactory(context.Background())
	if err != nil {
		return err
	}

	log.Print("Saving session")

	saver := saga.Repository().Save(sagaSession)
	err = uow.AddWorkUnit(saver)
	if err != nil {
		return err
	}

	log.Print("Inv")

	if firstStep.IsInvocable() {
		err = o.invokeStep(sagaSession, firstStep, uow)
		if err != nil {
			return err
		}
	} else {
		err = o.stepForwardAndInvoke(sagaSession, firstStep, sagaDef, uow)
		if err != nil {
			return err
		}
	}

	log.Print("Committing")

	err = uow.Commit()

	log.Print("Committed")
	return err
}

func (o *orchestrator[Tx]) invokeStep(session Session, curStep Step, uow *UnitOfWork[Tx]) error {
	var cmd Executable[Tx]
	var err error

	switch curStep.(type) {
	case remoteStep[Tx]:
		// Invoke the remote step.
		cmd = curStep.(remoteStep[Tx]).invocation(session)
	case localStep[Tx]:
		// Invoke the local step.
		cmd, err = curStep.(localStep[Tx]).invocation(session)
		if err != nil {
			return err
		}
	default:
		panic("unknown step type")
	}

	err = uow.AddWorkUnit(cmd)
	if err != nil {
		return err
	}

	session.SetPending(true)

	return nil
}

func (o *orchestrator[Tx]) stepForwardAndInvoke(session Session, curStep Step, def Definition, uow *UnitOfWork[Tx]) error {
	var err error

	log.Print("Setting state to is pending")

	nextStep := def.NextStep(curStep)
	if nextStep == nil {
		log.Print("Setting state to completed")
		session.SetState(StateCompleted)
		return nil
	}

	log.Print("Updating current step")

	err = session.UpdateCurrentStep(nextStep)
	if err != nil {
		return err
	}

	if nextStep.IsInvocable() {
		err = o.invokeStep(session, nextStep, uow)
		if err != nil {
			return err
		}

		log.Print("Done")

		return nil
	}

	err = o.stepForwardAndInvoke(session, nextStep, def, uow)
	if err != nil {
		return err
	}

	return nil
}

func (o *orchestrator[Tx]) stepBackwardAndCompensate(session Session, curStep Step, def Definition, uow *UnitOfWork[Tx]) error {
	var err error

	prevStep := def.PrevStep(curStep)
	if prevStep == nil {
		session.SetState(StateFailed)
		return nil
	}

	err = session.UpdateCurrentStep(prevStep)
	if err != nil {
		return err
	}

	if prevStep.IsCompensable() {
		session.SetState(StateIsCompensating)
		err = o.compensateStep(session, prevStep, uow)
		if err != nil {
			return err
		}

		return nil
	}

	err = o.stepBackwardAndCompensate(session, prevStep, def, uow)
	if err != nil {
		return err
	}

	return nil
}

func (o *orchestrator[Tx]) Orchestrate(saga Saga[Session, Tx], packet messagePacket) error {
	var uow *UnitOfWork[Tx]
	var err error

	origin := packet.Origin()
	if origin == "" {
		return errors.New("message origin is empty")
	}

	sagaSession, err := saga.Repository().Load(packet.Payload().SessionID())
	if err != nil {
		return err
	}

	if sagaSession.State() == StateCompleted || sagaSession.State() == StateFailed {
		return errors.New("session is already completed or failed")
	}

	currentStep := sagaSession.CurrentStep()
	if saga.Definition().Exists(currentStep) == false {
		return errors.New("session has no current step")
	}

	uow, err = o.uowFactory(context.Background())
	if err != nil {
		return err
	}

	if sagaSession.State() != StateIsCompensating &&
		sagaSession.State() != StateCompleted &&
		sagaSession.State() != StateFailed {
		// If the session is not compensating, completed, or failed, then it is in forward direction.
		err = o.handleInvocationResponse(sagaSession, origin, packet.Payload(), currentStep, saga.Definition(), uow)
	} else if sagaSession.State() == StateIsCompensating {
		// If the session is compensating, then it is in backward direction.
		err = o.handleCompensationResponse(sagaSession, origin, packet.Payload(), currentStep, saga.Definition(), uow)
	}

	if err != nil {
		return err
	}

	log.Print("Saving session")

	saver := saga.Repository().Save(sagaSession)
	err = uow.AddWorkUnit(saver)
	if err != nil {
		return err
	}

	err = uow.Commit()
	if err != nil {
		return err
	}

	log.Print("Committed orchestration")

	return nil
}

func (o *orchestrator[Tx]) handleInvocationResponse(session Session, origin ChannelName, msg Message, curStep Step, def Definition, uow *UnitOfWork[Tx]) error {
	session.SetPending(false)

	isFailure, err := o.isFailureInvocationResponse(origin, curStep)
	if err != nil {
		return err
	}

	if isFailure {
		var err error
		if curStep.MustBeCompleted() {
			err = o.retryInvocation(session, curStep, uow)
			return err
		}

		err = o.stepBackwardAndCompensate(session, curStep, def, uow)
		return err
	}

	err = o.stepForwardAndInvoke(session, curStep, def, uow)
	return err
}

func (o *orchestrator[Tx]) retryInvocation(session Session, step Step, uow *UnitOfWork[Tx]) error {
	if !step.MustBeCompleted() {
		return errors.New("step must be completed")
	}

	session.SetState(StateIsRetrying)
	err := o.invokeStep(session, step, uow)
	return err
}

func (o *orchestrator[Tx]) isFailureInvocationResponse(origin ChannelName, step Step) (bool, error) {
	var success ChannelName
	var failure ChannelName

	switch step.(type) {
	case remoteStep[Tx]:
		success = step.(remoteStep[Tx]).invokeEndpoint.SuccessResChannel()
		failure = step.(remoteStep[Tx]).invokeEndpoint.FailureResChannel()
	case localStep[Tx]:
		success = step.(localStep[Tx]).invokeEndpoint.SuccessResChannel()
		failure = step.(localStep[Tx]).invokeEndpoint.FailureResChannel()
	default:
		panic("unknown step type")
	}

	if failure != "" && origin == failure {
		return true, nil
	}

	if success != "" && origin == success {
		return false, nil
	}

	return false, errors.New("unknown chanaaanel")
}

func (o *orchestrator[Tx]) isFailureCompensationResponse(origin ChannelName, step Step) (bool, error) {
	var success ChannelName
	var failure ChannelName

	switch step.(type) {
	case remoteStep[Tx]:
		success = step.(remoteStep[Tx]).compEndpoint.SuccessResChannel()
		failure = step.(remoteStep[Tx]).compEndpoint.FailureResChannel()
	case localStep[Tx]:
		success = step.(localStep[Tx]).compEndpoint.SuccessResChannel()
		failure = step.(localStep[Tx]).compEndpoint.FailureResChannel()
	default:
		panic("unknown step type")
	}

	if failure != "" && origin == failure {
		return true, nil
	}

	if success != "" && origin == success {
		return false, nil
	}

	return false, errors.New("unknown channel")
}

func (o *orchestrator[Tx]) handleCompensationResponse(session Session, origin ChannelName, msg Message, curStep Step, def Definition, uow *UnitOfWork[Tx]) error {
	session.SetPending(false)

	isFailure, err := o.isFailureCompensationResponse(origin, curStep)
	if err != nil {
		return err
	}

	if isFailure {
		err = o.retryCompensation(session, curStep, uow)
		return err
	}

	err = o.stepBackwardAndCompensate(session, curStep, def, uow)
	return err
}

func (o *orchestrator[Tx]) retryCompensation(session Session, step Step, uow *UnitOfWork[Tx]) error {
	session.SetState(StateIsCompensating)
	err := o.compensateStep(session, step, uow)
	return err
}

func (o *orchestrator[Tx]) compensateStep(session Session, step Step, uow *UnitOfWork[Tx]) error {
	var cmd Executable[Tx]
	var err error

	switch step.(type) {
	case remoteStep[Tx]:
		msg := step.(remoteStep[Tx]).compEndpoint.SuccessResponseConstructor()(session)
		cmd = step.(remoteStep[Tx]).compEndpoint.CommandRepository().SaveMessage(msg)
	case localStep[Tx]:
		cmd, err = step.(localStep[Tx]).compEndpoint.handler(session)
		if err != nil {
			return err
		}
	default:
		panic("unknown step type")
	}

	err = uow.AddWorkUnit(cmd)
	if err != nil {
		return err
	}

	session.SetPending(true)

	return nil
}
