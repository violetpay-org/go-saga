package saga

import "errors"

var (
	ErrChannelAlreadyRegistered         = errors.New("channel already registered")
	ErrUnitOfWorkImmutable              = errors.New("unit of work is immutable because it has already been committed")
	ErrSessionCreationFailed            = errors.New("session is nil when creating a new session")
	ErrSessionIDEmpty                   = errors.New("session ID is empty")
	ErrSagaHasNoSteps                   = errors.New("saga has no steps")
	ErrUnknownMessageOrigin             = errors.New("message consumed but origin channel is unknown")
	ErrDeadSession                      = errors.New("session is already completed or failed")
	ErrSessionStepAndDefinitionMismatch = errors.New("session step and definition mismatch")
	ErrRetryCalledOnNonRetryingStep     = errors.New("retry called but step must be completed false")
	ErrRegisterInvalidSaga              = errors.New("saga is invalid, but tried to register")
	ErrSagaNotFound                     = errors.New("saga not found")
	ErrInvalidSagaStart                 = errors.New("start saga called with invalid parameters")
)
