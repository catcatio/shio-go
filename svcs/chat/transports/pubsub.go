package transports

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
)

type defaultPubsubHandler struct {
	handlers endpoints.PubsubMessageHandlers
}

func (d *defaultPubsubHandler) Serve(topic string, ctx context.Context, m pubsub.RawPubsubMessage) error {
	if handler, ok := d.handlers[topic]; ok {
		return handler(ctx, m)
	}

	return fmt.Errorf("handler not found for topic %s", topic)
}

func NewPubsubServer(handlers endpoints.PubsubMessageHandlers) pubsub.Handler {
	return &defaultPubsubHandler{
		handlers: handlers,
	}
}
