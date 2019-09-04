package repositories

import (
	"context"
	"errors"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	entities2 "github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/repositories/outgoing_handler"
)

var (
	ErrUnsupportedOutgoingEventType = errors.New("event type not support")
)

type OutgoingRepository interface {
	Handler(ctx context.Context, event *entities.OutgoingEvent) error
}

type outgoingRepo struct {
}

func NewOutgoingRepo() OutgoingRepository {
	return &outgoingRepo{}
}

func (o *outgoingRepo) Handler(ctx context.Context, event *entities.OutgoingEvent) error {
	var channelConfig entities2.ChannelConfig
	if config, ok := ctx.Value("channel_config").(entities2.ChannelConfig); !ok {
		return ErrChannelConfigNotSet
	} else {
		channelConfig = config
	}

	switch event.Type {
	case entities.OutgoingEventTypeMessage:
		// assume line for now
		client := outgoing_handler.NewLineClient(
			channelConfig.LineChatOptions.ChannelSecret,
			channelConfig.LineChatOptions.ChannelAccessToken,
		)
		err := client.TryReplyRaw(ctx, event.OutgoingMessage.ReplyToken, event.OutgoingMessage.RecipientID, event.OutgoingMessage.PayloadJson)
		return err

	default:
		return ErrUnsupportedOutgoingEventType
	}
}
