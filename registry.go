package saga

import (
	"errors"
	"sync"
)

func NewRegistry[Tx TxContext](orchestrator Orchestrator[Tx]) *Registry[Tx] {
	return &Registry[Tx]{
		sagas:        make([]Saga[Session, Tx], 0),
		mutex:        sync.Mutex{},
		orchestrator: orchestrator,
	}
}

func RegisterSagaTo[S Session, Tx TxContext](r *Registry[Tx], s Saga[S, Tx]) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.sagas = append(r.sagas, convertSaga(s))
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
			return err
		}
	}

	return nil
}

func (r *Registry[Tx]) StartSaga(sagaName string, sessionArgs map[string]interface{}) error {
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

	return r.orchestrator.StartSaga(*target, sessionArgs)
}
