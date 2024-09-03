package saga

import (
	"errors"
	"log"
	"sync"
)

func NewRegistry[Tx TxContext](orchestrator Orchestrator[Tx]) *Registry[Tx] {
	return &Registry[Tx]{
		sagas:        make([]Saga[Session, Tx], 0),
		mutex:        sync.Mutex{},
		orchestrator: orchestrator,
	}
}

func RegisterSagaTo[S Session, Tx TxContext](r *Registry[Tx], s Saga[S, Tx]) error {
	if s.Name() == "" {
		return errors.New("saga name is empty")
	}

	if r.HasSaga(s.Name()) {
		return errors.New("saga already registered")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.sagas = append(r.sagas, convertSaga(s))

	return nil
}

type Registry[Tx TxContext] struct {
	sagas        []Saga[Session, Tx]
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
			log.Print("Error in orchestrating: ", err)
			return err
		}
	}

	return nil
}

func (r *Registry[Tx]) StartSaga(sagaName string, sessionArgs map[string]interface{}) error {
	if sagaName == "" {
		return errors.New("saga name is empty")
	}

	if sessionArgs == nil {
		return errors.New("session args is nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	var target *Saga[Session, Tx]
	for _, s := range r.sagas {
		if s.Name() == sagaName {
			target = &s
			break
		}
	}

	if target == nil {
		return errors.New("saga not found")
	}

	log.Print("Starting saga: ", sagaName)

	return r.orchestrator.StartSaga(*target, sessionArgs)
}

func (r *Registry[Tx]) HasSaga(sagaName string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, s := range r.sagas {
		if s.Name() == sagaName {
			return true
		}
	}

	return false
}
