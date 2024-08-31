package saga

type State int

const (
	StateCommon State = iota
	StateCompleted
	StateFailed
	StateIsCompensating
	StateIsRetrying
)

type SessionFactory func(map[string]interface{}) Session

type Session interface {
	// ID returns the ID of the session.
	ID() string

	// CurrentStep returns the current step of the session.
	CurrentStep() string

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

type SessionRepository[Tx TxContext] interface {
	// Load finds a session by its ID.
	Load(id string) (Session, error)

	// Save saves a session.
	Save(sess Session) Executable[Tx]

	// Delete deletes a session.
	Delete(sess Session) Executable[Tx]
}
