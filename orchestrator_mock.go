package saga

type mockOrchestrator[Tx TxContext] struct {
}

func (m *mockOrchestrator[Tx]) Orchestrate(saga Saga[Session, Tx], packet messagePacket) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrchestrator[Tx]) StartSaga(saga Saga[Session, Tx], sessionArgs map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}
