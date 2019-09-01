package transports

import (
	"context"
	"errors"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
)

type defaultPubsubHandler struct {
	handlers endpoints.PubsubMessageHandlers
}

func (d *defaultPubsubHandler) Serve(topic string, ctx context.Context, m pubsub.RawPubsubMessage) error {
	if handler, ok := d.handlers[topic]; ok {
		return handler(ctx, m)
	}

	return ErrHandlerNotFound
}

func NewPubsubServer(handlers endpoints.PubsubMessageHandlers) pubsub.Handler {
	return &defaultPubsubHandler{
		handlers: handlers,
	}
}
