package usecases

import (
	"context"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
)

type IncomingEventUsecase interface {
	HandleEvents(ctx context.Context, in *entities.WebhookInput) (err error)
}

type incomingEventUsecase struct {
	channelOptionsRepo repositories.ChannelOptionsRepository
	incomingRepo       repositories.IncomingEventRepository
	log                *logger.Logger
}

func NewIncomingEventUsecase(channelOptionsRepo repositories.ChannelOptionsRepository, incomingRepo repositories.IncomingEventRepository) IncomingEventUsecase {
	return &incomingEventUsecase{
		channelOptionsRepo: channelOptionsRepo,
		incomingRepo:       incomingRepo,
		log:                logger.New("IncomingEvent"),
	}
}

func (i *incomingEventUsecase) HandleEvents(ctx context.Context, in *entities.WebhookInput) (err error) {
	log := i.log.WithServiceInfo("HandleIncomingEvents").WithRequestID(foundation.GetRequestIDFromContext(ctx))
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
