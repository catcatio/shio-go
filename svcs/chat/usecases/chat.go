package usecases

import (
	"context"
	entities2 "github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
)

type Chat interface {
	HandleIncomingEvents(ctx context.Context, in *entities.WebhookInput) error
	GetChannelConfig(ctx context.Context, channelID string) (*entities.ChannelConfig, error)
}

type chat struct {
	incomingRepo       repositories.IncomingEventRepository
	channelOptionsRepo repositories.ChannelOptionsRepository

	log *logger.Logger
}

func NewChat(incomingRepo repositories.IncomingEventRepository, channelOptionsRepo repositories.ChannelOptionsRepository) Chat {
	return &chat{
		incomingRepo:       incomingRepo,
		channelOptionsRepo: channelOptionsRepo,
		log:                logger.New("LineUsecase"),
	}
}

func (l *chat) HandleIncomingEvents(ctx context.Context, in *entities.WebhookInput) (err error) {
	log := l.log.WithServiceInfo("HandleIncomingEvents").WithRequestID(foundation.GetRequestIDFromContext(ctx))
	log.Infof("%d incoming event(s) from %s", len(in.Events), in.Provider)

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
