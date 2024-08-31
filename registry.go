package saga

import (
	"errors"
	"strings"
	"sync"
)

type Saga[Tx TxContext] struct {
	name       string
	definition Definition
	factory    SessionFactory
	repository SessionRepository[Tx]
}

func NewSaga[Tx TxContext](name string, def Definition, factory SessionFactory) Saga[Tx] {
	return Saga[Tx]{
		name:       name,
		definition: def,
		factory:    factory,
	}
}

func (s *Saga[Tx]) Name() string {
	return s.name
}

func (s *Saga[Tx]) Definition() Definition {
	return s.definition
}

func (s *Saga[Tx]) Repository() SessionRepository[Tx] {
	return s.repository
}

func (s *Saga[Tx]) createSession(args map[string]interface{}) Session {
	return s.factory(args)
}

func extractSagaName(sessid string) string {
	return strings.Split(sessid, "-")[0]
}

func (s *Saga[Tx]) hasPublishedSaga(sessid string) bool {
	sagaName := extractSagaName(sessid)
	return s.name == sagaName
}

func NewRegistry[Tx TxContext](orchestrator Orchestrator[Tx]) *Registry[Tx] {
	return &Registry[Tx]{
		sagas:        make([]Saga[Tx], 0),
		mutex:        sync.Mutex{},
		orchestrator: orchestrator,
	}
}

type Registry[Tx TxContext] struct {
	sagas        []Saga[Tx]
	mutex        sync.Mutex
	orchestrator Orchestrator[Tx]
}

func (r *Registry[Tx]) consumeMessage(packet messagePacket) error {
	sessID := packet.Payload().SessionID()

	type Orchestrations func() error
	orchestrations := make([]Orchestrations, 0)

	func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		for _, s := range r.sagas {
			if s.hasPublishedSaga(sessID) {
				f := func() error {
					return r.orchestrator.Orchestrate(s, packet)
				}
				orchestrations = append(orchestrations, f)
			}
		}
	}()

	for _, f := range orchestrations {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry[Tx]) StartSaga(sagaName string, sessionArgs map[string]interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var target *Saga[Tx]
	for _, s := range r.sagas {
		if s.Name() == sagaName {
			target = &s
			break
		}
	}

	if target == nil {
		return errors.New("saga not found")
	}

	return r.orchestrator.StartSaga(*target, sessionArgs)
}
