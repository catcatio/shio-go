package repositories

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
)

type IncomingEventRepository interface {
	Dispatch(ctx context.Context, event entities.IncomingEvent) error
}

type incomingEventRepository struct {
	pubsubClient *pubsub.Client
}

func NewIncomingEventRepository(pubsubClient *pubsub.Client) IncomingEventRepository {
	return &incomingEventRepository{pubsubClient: pubsubClient}
}

func (incomingEventRepository) Dispatch(ctx context.Context, event entities.IncomingEvent) error {
	panic("implement me")
}
