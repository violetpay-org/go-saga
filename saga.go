package saga

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func convertSaga[S Session, Tx TxContext](src Saga[S, Tx]) Saga[Session, Tx] {
	var factory SessionFactory[Session]
	factory = func(args map[string]interface{}) Session {
		sess := src.factory(args)
		return sess
	}

	var repository SessionRepository[Session, Tx]
	repository = &sessionRepository[Tx]{
		load:   func(id string) (Session, error) { return src.repository.Load(id) },
		save:   func(sess Session) Executable[Tx] { return src.repository.Save(sess.(S)) },
		delete: func(sess Session) Executable[Tx] { return src.repository.Delete(sess.(S)) },
	}

	return NewSaga[Session, Tx](src.name, src.definition, factory, repository)
}

type Saga[S Session, Tx TxContext] struct {
	name       string
	definition Definition
	factory    SessionFactory[S]
	repository SessionRepository[S, Tx]
}

func NewSaga[S Session, Tx TxContext](name string, def Definition, factory SessionFactory[S], repository SessionRepository[S, Tx]) Saga[S, Tx] {
	return Saga[S, Tx]{
		name:       name,
		definition: def,
		factory:    factory,
		repository: repository,
	}
}

func (s *Saga[S, Tx]) Name() string {
	return s.name
}

func (s *Saga[S, Tx]) Definition() Definition {
	return s.definition
}

func (s *Saga[S, Tx]) Repository() SessionRepository[S, Tx] {
	return s.repository
}

func (s *Saga[S, Tx]) createSession(args map[string]interface{}) S {
	if args["id"] == nil {
		args["id"] = fmt.Sprintf("%s-%s", s.name, uuid.New().String())
	}

	return s.factory(args)
}

func extractSagaName(sessid string) string {
	return strings.Split(sessid, "-")[0]
}

func (s *Saga[S, Tx]) hasPublishedSaga(sessid string) bool {
	sagaName := extractSagaName(sessid)
	return s.name == sagaName
}
