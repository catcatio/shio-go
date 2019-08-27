package repositories

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
)

type SendMessageRepository interface {
	Send(ctx context.Context, input entities.SendMessageInput) error
}

type sendMessageRepository struct {
	pubsubClient *pubsub.Clients
}

func NewSendMessageRepository(c *pubsub.Clients) SendMessageRepository {
	return &sendMessageRepository{
		pubsubClient: c,
	}
}

func (s *sendMessageRepository) Send(ctx context.Context, input entities.SendMessageInput) error {
	return s.pubsubClient.PublishSendMessageInput(ctx, input)
}
