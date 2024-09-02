package saga

import "errors"

type mockSession struct {
	id          string
	currentStep Step
	pending     bool
	state       State
	mockField   string
}

func (m *mockSession) ID() string {
	return m.id
}

func (m *mockSession) CurrentStep() Step {
	return m.currentStep
}

func (m *mockSession) UpdateCurrentStep(step Step) error {
	m.currentStep = step
	return nil
}

func (m *mockSession) IsPending() bool {
	return m.pending
}

func (m *mockSession) SetPending(pending bool) {
	m.pending = pending
}

func (m *mockSession) State() State {
	return m.state
}

func (m *mockSession) SetState(state State) {
	m.state = state
}

func newMockSessionRepository[S Session, Tx TxContext]() *mockSessionRepository[S, Tx] {
	return &mockSessionRepository[S, Tx]{
		sessions: make(map[string]*S),
	}
}

type mockSessionRepository[S Session, Tx TxContext] struct {
	sessions map[string]*S
}

func (m *mockSessionRepository[S, Tx]) Load(id string) (s S, err error) {
	if sess, ok := m.sessions[id]; ok {
		return *sess, nil
	}

	return s, errors.New("session not found")
}

func (m *mockSessionRepository[S, Tx]) Save(sess S) Executable[Tx] {
	return func(ctx Tx) error {
		m.sessions[sess.ID()] = &sess
		return nil
	}
}

func (m *mockSessionRepository[S, Tx]) Delete(sess S) Executable[Tx] {
	return func(ctx Tx) error {
		delete(m.sessions, sess.ID())
		return nil
	}
}
