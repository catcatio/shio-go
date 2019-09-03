package usecases

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/octofoxio/foundation/logger"
)

type OutgoingEventUsecase interface {
	Handle(ctx context.Context, input *entities.OutgoingEvent) error
}

type outgoingEventUsecase struct {
	channelConfigRepo repositories.ChannelConfigRepository
	pubsubRepo        repositories.PubsubChannelRepository
	log               *logger.Logger
}

func (s *outgoingEventUsecase) Handle(ctx context.Context, input *entities.OutgoingEvent) error {
	log := s.log.WithServiceInfo("Handle")
	channelConfig, err := s.channelConfigRepo.Get(ctx, input.ChannelID)

	if err != nil {
		log.WithError(err).WithField("channelID", input.ChannelID).Error("unable to get channel options")
		return err
	}

	// assume event message and line
	bot, _ := linebot.New(channelConfig.LineChatOptions.ChannelSecret, channelConfig.LineChatOptions.ChannelAccessToken)
	_, err = bot.ReplyMessage(input.OutgoingMessage.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%v", input.OutgoingMessage.Payload))).Do()

	if err != nil {
		return err
	}

	return nil

}

func NewOutgoingEventUsecase(channelConfigRepo repositories.ChannelConfigRepository, pubsubRepo repositories.PubsubChannelRepository) OutgoingEventUsecase {
	return &outgoingEventUsecase{
		channelConfigRepo: channelConfigRepo,
		pubsubRepo:        pubsubRepo,
		log:               logger.New("OutgoingEventUsecase"),
	}
}
