package messageRelayer

import (
	"errors"
	"github.com/violetpay-org/go-saga"
	"sync"
)

type ChannelRegistry[Tx saga.TxContext] interface {
	RegisterChannel(channel Channel[Tx]) error

	// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
	Range(func(name saga.ChannelName, channel Channel[Tx]) bool)

	// Find returns the channel with the given name.
	Find(name saga.ChannelName) Channel[Tx]

	// Has returns true if the channel with the given name exists.
	Has(name saga.ChannelName) bool
}

// Channel is a channel that can send AbstractMessage to *Somewhere*.
// Relayer will get messages from outbox repository (messageRepository), and uses this channel to send message to the other remote service.
type Channel[Tx saga.TxContext] interface {
	saga.AbstractChannel[Tx]
}

func NewChannel[M saga.Message, Tx saga.TxContext](name saga.ChannelName, registry *saga.Registry[Tx], repository saga.AbstractMessageRepository[M, Tx], send func(message saga.Message) error) Channel[Tx] {
	return &channel[Tx]{name: name, registry: registry, repository: saga.ConvertMessageRepository(repository), send: send}
}

type channel[Tx saga.TxContext] struct {
	name       saga.ChannelName
	registry   *saga.Registry[Tx]
	repository saga.AbstractMessageRepository[saga.Message, Tx]
	send       func(message saga.Message) error
}

func (c *channel[Tx]) Name() saga.ChannelName {
	return c.name
}

func (c *channel[Tx]) Send(message saga.Message) error {
	return c.send(message)
}

func (c *channel[Tx]) Repository() saga.AbstractMessageRepository[saga.Message, Tx] {
	return c.repository
}

type channelRegistry[Tx saga.TxContext] struct {
	channels sync.Map
}

func NewChannelRegistry[Tx saga.TxContext]() ChannelRegistry[Tx] {
	return &channelRegistry[Tx]{}
}

func (r *channelRegistry[Tx]) RegisterChannel(channel Channel[Tx]) error {
	if _, ok := r.channels.Load(channel.Name()); ok {
		return errors.New("channel already exists")
	}
	r.channels.Store(channel.Name(), channel)
	return nil
}

func (r *channelRegistry[Tx]) Range(f func(name saga.ChannelName, channel Channel[Tx]) bool) {
	r.channels.Range(func(key, value interface{}) bool {
		return f(key.(saga.ChannelName), value.(Channel[Tx]))
	})
}

func (r *channelRegistry[Tx]) Find(name saga.ChannelName) Channel[Tx] {
	if channel, ok := r.channels.Load(name); ok {
		return channel.(Channel[Tx])
	}
	return nil
}

func (r *channelRegistry[Tx]) Has(name saga.ChannelName) bool {
	_, ok := r.channels.Load(name)
	return ok
}
