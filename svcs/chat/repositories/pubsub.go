package repositories

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
)

type PubsubChannelRepository interface {
	PublishIncomingEvent(ctx context.Context, input *entities.IncomingEvent) error
	PublishOutgoingEventInput(ctx context.Context, input *entities.OutgoingEvent) error
	PublishFulfillmentInput(ctx context.Context, input *entities.IncomingEvent) error
}

type pubsubChannelRepository struct {
	pubsubClient *pubsub.Clients
}

func NewPubsubChannelRepository(pubsubClient *pubsub.Clients) PubsubChannelRepository {
	return &pubsubChannelRepository{pubsubClient: pubsubClient}
}

func (p *pubsubChannelRepository) PublishIncomingEvent(ctx context.Context, input *entities.IncomingEvent) error {
	return p.pubsubClient.PublishIncomingEvent(ctx, input)
}

func (p *pubsubChannelRepository) PublishOutgoingEventInput(ctx context.Context, input *entities.OutgoingEvent) error {
	return p.pubsubClient.PublishOutgoingEventInput(ctx, input)
}

func (p *pubsubChannelRepository) PublishFulfillmentInput(ctx context.Context, input *entities.IncomingEvent) error {
	return p.pubsubClient.PublishFulfillmentInput(ctx, input)
}
