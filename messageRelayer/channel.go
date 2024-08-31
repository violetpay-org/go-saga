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

type Channel[Tx saga.TxContext] interface {
	saga.AbstractChannel[Tx]
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
