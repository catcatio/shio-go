package usecases

import (
	"context"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
)

type IncomingEventUsecase interface {
	GetChannelConfig(ctx context.Context, channelID string) (*entities.ChannelConfig, error)
	HandleEvents(ctx context.Context, in *entities.WebhookInput) (err error)
}

type incomingEventUsecase struct {
	channelConfigRepo repositories.ChannelConfigRepository
	incomingRepo      repositories.IncomingEventRepository
	log               *logger.Logger
}

func NewIncomingEventUsecase(channelConfigRepo repositories.ChannelConfigRepository, incomingRepo repositories.IncomingEventRepository) IncomingEventUsecase {
	return &incomingEventUsecase{
		channelConfigRepo: channelConfigRepo,
		incomingRepo:      incomingRepo,
		log:               logger.New("IncomingEvent"),
	}
}

func (i *incomingEventUsecase) HandleEvents(ctx context.Context, in *entities.WebhookInput) (err error) {
	log := i.log.WithServiceInfo("HandleEvents").WithRequestID(foundation.GetRequestIDFromContext(ctx))
	log.Infof("%d incoming event(s) from %s", len(in.Events), in.Provider)

	// forward event
	for _, e := range in.Events {
		err = i.incomingRepo.Dispatch(ctx, e)
		if err != nil {
			log.WithError(err).Error("dispatch event failed")
			return
		}
	}

	return
}

func (i *incomingEventUsecase) GetChannelConfig(ctx context.Context, channelID string) (*entities.ChannelConfig, error) {
	return i.channelConfigRepo.Get(ctx, channelID)
}
