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

func (m *mockAbstractMessageRepository) GetMessagesFromOutbox(batchSize int) ([]Message, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) GetMessagesFromDeadLetter(batchSize int) ([]Message, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveMessage(message Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveMessages(messages []Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveDeadLetter(message Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) SaveDeadLetters(message []Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteMessage(message Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteMessages(messages []Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteDeadLetter(message Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}

func (m *mockAbstractMessageRepository) DeleteDeadLetters(message Message) Executable[mockTxContext] {
	//TODO implement me
	panic("implement me")
}
