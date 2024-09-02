package saga

type State int

const (
	StateCommon State = iota
	StateCompleted
	StateFailed
	StateIsCompensating
	StateIsRetrying
)

type SessionFactory[S Session] func(map[string]interface{}) S

//type SessionID string

type Session interface {
	// ID returns the ID of the session.
	ID() string

	// CurrentStep returns the current step of the session.
	CurrentStep() Step

	// UpdateCurrentStep updates the current step of the session.
	UpdateCurrentStep(step Step) error

	// IsPending returns true if the session is pending. Session can be pending state when following conditions are met:
	//
	// - The session is in.
	//
	// - The session is in.
	//
	// - The session is in.
	IsPending() bool

	// SetPending sets the pending state of the session.
	SetPending(pending bool)

	// State returns the current state of the session.
	State() State

	// SetState sets the state of the session.
	SetState(state State)
}

type SessionRepository[S Session, Tx TxContext] interface {
	// Load finds a session by its ID.
	Load(id string) (S, error)

	// Save saves a session.
	Save(sess S) Executable[Tx]

	// Delete deletes a session.
	Delete(sess S) Executable[Tx]
}

type sessionRepository[Tx TxContext] struct {
	load   func(id string) (Session, error)
	save   func(sess Session) Executable[Tx]
	delete func(sess Session) Executable[Tx]
}

func (s *sessionRepository[Tx]) Load(id string) (Session, error) {
	return s.load(id)
}

func (s *sessionRepository[Tx]) Save(sess Session) Executable[Tx] {
	return s.save(sess)
}

func (s *sessionRepository[Tx]) Delete(sess Session) Executable[Tx] {
	return s.delete(sess)
}
