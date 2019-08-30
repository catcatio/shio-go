package repositories

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
)

type IncomingEventRepository interface {
	Dispatch(ctx context.Context, event *entities.IncomingEvent) error
}

type incomingEventRepository struct {
	pubsubClient *pubsub.Clients
}

func NewIncomingEventRepository(pubsubClient *pubsub.Clients) IncomingEventRepository {
	return &incomingEventRepository{pubsubClient: pubsubClient}
}

func (i *incomingEventRepository) Dispatch(ctx context.Context, event *entities.IncomingEvent) error {
	return i.pubsubClient.PublishIncomingEvent(ctx, event)
}
