package saga

import "time"

func newMockMessage() mockMessage {
	return mockMessage{}
}

type mockMessage struct {
}

func (m mockMessage) ID() string {
	//TODO implement me
	panic("implement me")
}

func (m mockMessage) SessionID() string {
	//TODO implement me
	panic("implement me")
}

func (m mockMessage) Trigger() string {
	//TODO implement me
	panic("implement me")
}

func (m mockMessage) CreatedAt() time.Time {
	//TODO implement me
	panic("implement me")
}

func (m mockMessage) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func newMockAbstractMessageRepository() *mockAbstractMessageRepository {
	return &mockAbstractMessageRepository{}
}

type mockAbstractMessageRepository struct {
}

func (m *mockAbstractMessageRepository) GetMessagesFromOutbox(batchSize int) ([]mockMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) GetMessagesFromDeadLetter(batchSize int) ([]mockMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveMessage(message mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveMessages(messages []mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveDeadLetter(message mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveDeadLetters(message []mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteMessage(message mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteMessages(messages []mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteDeadLetter(message mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteDeadLetters(messages []mockMessage) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}
