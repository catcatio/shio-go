package pubsub

import "context"

type RawPubsubMessage struct {
	Data []byte `json:"data"`
}

type Handler interface {
	Serve(topic string, ctx context.Context, m RawPubsubMessage) error
}
