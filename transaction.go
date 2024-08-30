package saga

import (
	"context"
	"errors"
)

//type MySQLTxContext struct {
//	Tx *sql.Tx
//}
//
//type MySQLTxHandler struct {
//	DB *sql.DB
//}
//
//func (ctor *MySQLTxHandler) BeginTx(context context.Context) (ctx MySQLTxContext, error error) {
//	tx, err := ctor.DB.BeginTx(context, nil)
//	if err != nil {
//		return MySQLTxContext{}, err
//	}
//
//	return MySQLTxContext{Tx: tx}, nil
//}
//
//func (ctor *MySQLTxHandler) Commit(ctx MySQLTxContext) error {
//	if err := ctx.Tx.Commit(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (ctor *MySQLTxHandler) Rollback(ctx MySQLTxContext) error {
//	if err := ctx.Tx.Rollback(); err != nil {
//		return err
//	}
//
//	return nil
//}

type TxContext interface{}

type TxHandler[Tx TxContext] interface {
	BeginTx(ctx context.Context) (tx Tx, error error)
	Commit(ctx Tx) error
	Rollback(ctx Tx) error
}

type Executable[Tx TxContext] func(ctx Tx) error

func CombineExecutables[Tx TxContext](executables ...Executable[Tx]) Executable[Tx] {
	return func(ctx Tx) error {
		for _, executable := range executables {
			err := executable(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

type UnitOfWorkFactory[Tx TxContext] func(ctx context.Context) (UnitOfWork[Tx], error)

type UnitOfWork[Tx TxContext] struct {
	handler  TxHandler[Tx]
	unitChan chan Executable[Tx]
	ctx      context.Context
	commited bool
}

func NewUnitOfWork[Tx TxContext](ctx context.Context, handler TxHandler[Tx]) *UnitOfWork[Tx] {
	return &UnitOfWork[Tx]{
		handler:  handler,
		unitChan: make(chan Executable[Tx]),
		ctx:      ctx,
	}
}

func (u *UnitOfWork[Tx]) AddWorkUnit(workUnit Executable[Tx]) error {
	if u.commited {
		return errors.New("cannot add work unit to immutable unit of work")
	}

	u.unitChan <- workUnit
	return nil
}

func (u *UnitOfWork[Tx]) commitExecutables(ctx Tx) error {
	if u.commited {
		return errors.New("cannot commit immutable unit of work, already commited")
	}

	u.commited = true

	backupChan := make(chan Executable[Tx])
	errors := make(chan error)
	defer close(errors)

commitLoop:
	for {
		select {
		case executable := <-u.unitChan:
			backupChan <- executable
			err := executable(ctx)
			if err != nil {
				errors <- err
			}
		default:
			close(u.unitChan)
			break commitLoop
		}
	}

	select {
	case err := <-errors:
		u.unitChan = backupChan
		return err
	default:
		break
	}

	close(backupChan)
	return nil
}

func (u *UnitOfWork[Tx]) Commit() error {
	tx, err := u.handler.BeginTx(u.ctx)
	defer u.handler.Rollback(tx)
	if err != nil {
		return err
	}

	err = u.commitExecutables(tx)
	if err != nil {
		return err
	}

	err = u.handler.Commit(tx)
	if err != nil {
		return err
	}

	return nil
}
