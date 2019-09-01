package repositories

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
)

type OutgoingEventRepository interface {
	Send(ctx context.Context, input *entities.OutgoingEvent) error
}

type outgoingEventRepository struct {
	pubsubClient *pubsub.Clients
}

func NewOutgoingEventRepository(c *pubsub.Clients) OutgoingEventRepository {
	return &outgoingEventRepository{
		pubsubClient: c,
	}
}

func (s *outgoingEventRepository) Send(ctx context.Context, input *entities.OutgoingEvent) error {
	return s.pubsubClient.PublishOutgoingEventInput(ctx, input)
}
