package usecases

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/octofoxio/foundation/logger"
)

type SendMessageUsecase interface {
	HandleMessage(ctx context.Context, input *entities.SendMessageInput) error
}

type sendMessageUsecase struct {
	sendMessageRepo   repositories.SendMessageRepository
	channelConfigRepo repositories.ChannelConfigRepository
	log               *logger.Logger
}

func NewSendMessageUsecase(channelConfigRepo repositories.ChannelConfigRepository, sendMessageRepo repositories.SendMessageRepository) SendMessageUsecase {
	return &sendMessageUsecase{
		channelConfigRepo: channelConfigRepo,
		sendMessageRepo:   sendMessageRepo,
		log:               logger.New("SendMessageUsecase"),
	}
}

func (s *sendMessageUsecase) HandleMessage(ctx context.Context, input *entities.SendMessageInput) error {
	log := s.log.WithServiceInfo("HandleMessage")
	channelConfig, err := s.channelConfigRepo.Get(ctx, input.ChannelID)

	if err != nil {
		log.WithError(err).WithField("channelID", input.ChannelID).Error("unable to get channel options")
		return err
	}

	// assume line for now
	// TODO get messageSender by provider
	bot, _ := linebot.New(channelConfig.LineChatOptions.ChannelSecret, channelConfig.LineChatOptions.ChannelAccessToken)

	_, err = bot.ReplyMessage(input.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%v", input.Payload))).Do()

	if err != nil {
		return err
	}

	return nil
}
