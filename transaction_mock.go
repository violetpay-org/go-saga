package saga

import (
	"context"
)

func newMockTxContext() *mockTxContext {
	return &mockTxContext{}
}

type mockTxContext struct{}

func newMockTxHandler() *mockTxHandler {
	return &mockTxHandler{}
}

type mockTxHandler struct{}

func (m *mockTxHandler) BeginTx(ctx context.Context) (tx mockTxContext, error error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxHandler) Commit(ctx mockTxContext) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxHandler) Rollback(ctx mockTxContext) error {
	//TODO implement me
	panic("implement me")
}
