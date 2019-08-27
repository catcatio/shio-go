package usecases

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/nub"
	entities2 "github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
)

type Chat interface {
	HandleIncomingEvents(ctx context.Context, in *entities.WebhookInput) error
	GetChannelConfig(ctx context.Context, channelID string) (*entities.ChannelConfig, error)
	SendMessage(ctx context.Context, input *entities2.SendMessageInput) error
}

type chat struct {
	incomingRepo       repositories.IncomingEventRepository
	channelOptionsRepo repositories.ChannelOptionsRepository
	sendMessageRepo    repositories.SendMessageRepository

	log *logger.Logger
}

func (l *chat) SendMessage(ctx context.Context, input *entities2.SendMessageInput) error {
	//channelOptions, err := chat.GetChannelConfig(ctx, input.ChannelID)
	//
	//if err != nil {
	//	l.log.WithError(err).WithField("channelID", input.ChannelID).Error("unable to get channel options")
	//	return err
	//}
	//
	//bot, _ := linebot.New(channelOptions.LineChatOptions.ChannelSecret, channelOptions.ChannelAccessToken)
	//bot.ReplyMessage()
	l.log.WithField("input", input).WithServiceInfo("SendMessage").Info("send a message")
	return nil
}

func NewChat(incomingRepo repositories.IncomingEventRepository, channelOptionsRepo repositories.ChannelOptionsRepository, sendMessageRepo repositories.SendMessageRepository) Chat {
	return &chat{
		incomingRepo:       incomingRepo,
		channelOptionsRepo: channelOptionsRepo,
		sendMessageRepo:    sendMessageRepo,
		log:                logger.New("LineUsecase"),
	}
}

func (l *chat) HandleIncomingEvents(ctx context.Context, in *entities.WebhookInput) (err error) {
	log := l.log.WithServiceInfo("HandleIncomingEvents").WithRequestID(foundation.GetRequestIDFromContext(ctx))
	log.Infof("%d incoming event(s) from %s", len(in.Events), in.Provider)

	for _, e := range in.Events {
		_ = l.sendMessageRepo.Send(ctx, entities2.SendMessageInput{
			RequestID:   nub.NewID(),
			ChannelID:   in.ChannelID,
			ReplyToken:  e.ReplyToken,
			Provider:    e.Provider,
			RecipientID: e.Source.GetSourceID(),
			Payload:     linebot.NewTextMessage(fmt.Sprintf("thank for your message, %v", e.Message)),
		})
	}

	// TODO: get event dispatcher config by provider
	// forward event
	for _, e := range in.Events {
		err := l.incomingRepo.Dispatch(ctx, entities2.IncomingEvent{})
		log.Println(e)
		if err != nil {
			log.WithError(err).Error("dispatch event failed")
			return err
		}
	}

	return nil
}

func (l *chat) GetChannelConfig(ctx context.Context, channelID string) (*entities.ChannelConfig, error) {
	return l.channelOptionsRepo.Get(ctx, channelID)
}
