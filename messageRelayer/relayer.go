package messageRelayer

import (
	"context"
	"errors"
	"github.com/violetpay-org/go-saga"
	"sync"
	"sync/atomic"
)

type Relayer[Tx saga.TxContext] struct {
	batchSize int
	mutex     sync.Mutex
	registry  ChannelRegistry[Tx]

	unitOfWork        *saga.UnitOfWork[Tx]
	unitOfWorkFactory saga.UnitOfWorkFactory[Tx]
}

func New[Tx saga.TxContext](
	batchSize int,
	registry ChannelRegistry[Tx],
	factory saga.UnitOfWorkFactory[Tx],
) BatchJob {
	return &Relayer[Tx]{
		batchSize:         batchSize,
		registry:          registry,
		unitOfWorkFactory: factory,
		mutex:             sync.Mutex{},
	}
}

func (r *Relayer[Tx]) Execute() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.createUnitOfWork(context.Background()); err != nil {
		return err
	}

	err := r.relayAndSave()
	if err != nil {
		return err
	}

	err = r.commitUnitOfWork()
	if err != nil {
		return err
	}

	return nil
}

func (r *Relayer[Tx]) createUnitOfWork(ctx context.Context) error {
	if r.unitOfWork != nil {
		return errors.New("duplicate unit of work create")
	}

	uow, err := r.unitOfWorkFactory(ctx)
	if err != nil {
		return err
	}

	r.unitOfWork = uow
	return nil
}

func (r *Relayer[Tx]) commitUnitOfWork() error {
	if r.unitOfWork == nil {
		return errors.New("failed to relay messages")
	}

	err := r.unitOfWork.Commit()
	r.unitOfWork = nil

	return err
}

func (r *Relayer[Tx]) relayAndSave() error {
	remaining := &atomic.Int64{}
	remaining.Store(int64(r.batchSize))

	published, failed := r.publishFromOutbox(remaining)
	defer published.close()
	defer failed.close()

	var publishedErr, failedErr error
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			name, messages, ok := published.popMessagesChannelPair()
			if !ok {
				return
			}

			err := r.deleteMessagesFromOutbox(name, messages)
			if err != nil {
				publishedErr = err
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			name, messages, ok := failed.popMessagesChannelPair()
			if !ok {
				return
			}

			err := r.saveDeadLetters(name, messages)
			if err != nil {
				failedErr = err
				return
			}
		}
	}()

	wg.Wait()

	if publishedErr != nil {
		return publishedErr
	}

	if failedErr != nil {
		return failedErr
	}

	published = nil
	failed = nil

	published, failed = r.publishFromDeadLetters(remaining)
	defer published.close()
	defer failed.close()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			name, messages, ok := published.popMessagesChannelPair()
			if !ok {
				return
			}

			err := r.deleteMessagesFromDeadLetters(name, messages)
			if err != nil {
				publishedErr = err
				return
			}
		}
	}()

	wg.Wait()

	if publishedErr != nil {
		return publishedErr
	}

	return nil
}

func (r *Relayer[Tx]) saveDeadLetters(name saga.ChannelName, messages <-chan saga.Message) error {
	channel := r.registry.Find(name)
	if channel == nil {
		return errors.New("channel not found")
	}

	repo := channel.Repository()

saveDeadLetters:
	for {
		select {
		case message := <-messages:
			deleteCmd := repo.DeleteMessage(message)
			saveCmd := repo.SaveDeadLetter(message)

			err := r.unitOfWork.AddWorkUnit(deleteCmd)
			if err != nil {
				return err
			}

			err = r.unitOfWork.AddWorkUnit(saveCmd)
			if err != nil {
				return err
			}
		default:
			break saveDeadLetters
		}
	}

	return nil
}

func (r *Relayer[Tx]) deleteMessagesFromOutbox(name saga.ChannelName, messages <-chan saga.Message) error {
	channel := r.registry.Find(name)
	if channel == nil {
		return errors.New("channel not found")
	}

	repo := channel.Repository()

deleteMessages:
	for {
		select {
		case message := <-messages:
			cmd := repo.DeleteMessage(message)
			err := r.unitOfWork.AddWorkUnit(cmd)
			if err != nil {
				return err
			}
		default:
			break deleteMessages
		}
	}

	return nil
}

func (r *Relayer[Tx]) deleteMessagesFromDeadLetters(name saga.ChannelName, messages <-chan saga.Message) error {
	channel := r.registry.Find(name)
	if channel == nil {
		return errors.New("channel not found")
	}

	repo := channel.Repository()

deleteDeadLetters:
	for {
		select {
		case message := <-messages:
			cmd := repo.DeleteDeadLetter(message)
			err := r.unitOfWork.AddWorkUnit(cmd)
			if err != nil {
				return err
			}
		default:
			break deleteDeadLetters
		}
	}

	return nil
}

func (r *Relayer[Tx]) publishFromOutbox(remaining *atomic.Int64) (published *messagesByChannel, failed *messagesByChannel) {
	messageFunc := func(repo saga.AbstractMessageLoadRepository, batchSize int) ([]saga.Message, error) {
		return repo.GetMessagesFromOutbox(batchSize)
	}

	published, failed = r.publish(remaining, messageFunc)
	return
}

func (r *Relayer[Tx]) publishFromDeadLetters(remaining *atomic.Int64) (published *messagesByChannel, failed *messagesByChannel) {
	messageFunc := func(repo saga.AbstractMessageLoadRepository, batchSize int) ([]saga.Message, error) {
		return repo.GetMessagesFromDeadLetter(batchSize)
	}

	published, failed = r.publish(remaining, messageFunc)
	return
}

func (r *Relayer[Tx]) publish(remaining *atomic.Int64, messageFunc func(repo saga.AbstractMessageLoadRepository, batchSize int) ([]saga.Message, error)) (published *messagesByChannel, failed *messagesByChannel) {
	published = newMessagesByChannel(r.batchSize)
	failed = newMessagesByChannel(r.batchSize)

	r.registry.Range(func(name saga.ChannelName, channel Channel[Tx]) bool {
		batchSize := int(remaining.Load())
		if batchSize <= 0 {
			return false
		}

		repo := channel.Repository()
		messages, err := messageFunc(repo, batchSize)
		if err != nil {
			return false
		}

		wg := sync.WaitGroup{}
		for _, message := range messages {
			wg.Add(1)
			go func(message saga.Message) {
				defer wg.Done()
				err := channel.Send(message)
				if err != nil {
					failed.pushMessage(name, message)
				} else {
					published.pushMessage(name, message)
				}

				remaining.Add(-1)
			}(message)
		}

		wg.Wait()

		return true
	})

	return
}

type messagesByChannel struct {
	once              sync.Once
	baseBatchSize     int
	messagesToChannel sync.Map
	channelNames      chan saga.ChannelName
}

func newMessagesByChannel(baseBatchSize int) *messagesByChannel {
	return &messagesByChannel{
		once:              sync.Once{},
		baseBatchSize:     baseBatchSize,
		channelNames:      make(chan saga.ChannelName, 100), // max 100 channels
		messagesToChannel: sync.Map{},
	}
}

func (m *messagesByChannel) pushMessage(channelName saga.ChannelName, message saga.Message) {
	channelMessages, loaded := m.messagesToChannel.LoadOrStore(channelName, make(chan saga.Message, m.baseBatchSize))
	if !loaded {
		m.channelNames <- channelName
	}

	channelMessages.(chan saga.Message) <- message
}

func (m *messagesByChannel) popMessagesChannelPair() (channelName saga.ChannelName, messages <-chan saga.Message, ok bool) {
	for {
		select {
		case channelName = <-m.channelNames:
			loaded, ok := m.messagesToChannel.Load(channelName)
			if !ok {
				return "", nil, false
			}

			messages, ok = loaded.(chan saga.Message)
			return channelName, messages, ok
		default:
			return
		}
	}
}

func (m *messagesByChannel) close() {
	m.once.Do(
		func() {
			m.messagesToChannel.Range(
				func(key, value interface{}) bool {
					close(value.(chan saga.Message))
					return true
				},
			)
		},
	)
}
