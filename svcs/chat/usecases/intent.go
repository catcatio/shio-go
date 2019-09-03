package usecases

import (
	"context"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation/logger"
)

type IntentUsecase interface {
	HandleEvents(ctx context.Context, in *entities.IncomingEvent) (err error)
}

type intentUsecase struct {
	channelConfigRepo repositories.ChannelConfigRepository
	intentRepo        repositories.IntentRepository
	pubsubRepo        repositories.PubsubChannelRepository
	log               *logger.Logger
}

func NewIntentUsecase(channelConfigRepo repositories.ChannelConfigRepository, intentRepo repositories.IntentRepository, pubsubRepo repositories.PubsubChannelRepository) IntentUsecase {
	return &intentUsecase{
		channelConfigRepo: channelConfigRepo,
		intentRepo:        intentRepo,
		pubsubRepo:        pubsubRepo,
		log:               logger.New("IntentUsecase"),
	}
}

func (i *intentUsecase) HandleEvents(ctx context.Context, in *entities.IncomingEvent) (err error) {
	log := i.log.WithServiceInfo("HandleEvents").WithRequestID(shio.ReqIDFromContext(ctx))
	log.Println(in)
	intent, err := i.intentRepo.Detect(ctx, in)

	if err != nil {
		log.WithError(err).Error("detect intent failed")
		return err
	}

	in.Intent = intent
	log.Println(intent)

	err = i.pubsubRepo.PublishFulfillmentInput(ctx, in)
	if err != nil {
		log.WithError(err).Error("publish fulfillment failed")
		return err
	}

	return nil
}
