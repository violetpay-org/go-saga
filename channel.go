package saga

type ChannelName string

type AbstractChannel[Tx TxContext] interface {
	// Name returns the name of the channel
	Name() ChannelName

	// Send sends a AbstractMessage to the channel
	Send(message Message) error

	Repository() AbstractMessageRepository[Tx]
}

// parseMessageToPacket parses the AbstractMessage to a AbstractMessage packet
func parseMessageToPacket[Tx TxContext](channel AbstractChannel[Tx], message Message) messagePacket {
	return newMessagePacket(channel.Name(), message)
}

// Channel is a channel that can send AbstractMessage to *Registry*.
type Channel[Tx TxContext] interface {
	AbstractChannel[Tx]
}

func NewChannel[Tx TxContext](name string, registry *Registry[Tx], repository AbstractMessageRepository[Tx]) Channel[Tx] {
	return &channel[Tx]{name: ChannelName(name), registry: registry, repository: repository}
}

type channel[Tx TxContext] struct {
	name       ChannelName
	registry   *Registry[Tx]
	repository AbstractMessageRepository[Tx]
}

func (c *channel[Tx]) Name() ChannelName {
	return c.name
}

func (c *channel[Tx]) Send(message Message) error {
	packet := parseMessageToPacket[Tx](c, message)
	return c.registry.consumeMessage(packet)
}

func (c *channel[Tx]) Repository() AbstractMessageRepository[Tx] {
	return c.repository
}
