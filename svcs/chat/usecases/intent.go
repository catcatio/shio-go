package usecases

import (
	"context"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	line2 "github.com/catcatio/shio-go/utils/line"
	"github.com/octofoxio/foundation/logger"
)

type IntentUsecase interface {
	HandleEvents(ctx context.Context, in *entities.IncomingEvent) (err error)
}

type intentUsecase struct {
	intentRepo repositories.IntentRepository
	pubsubRepo repositories.PubsubChannelRepository
	log        *logger.Logger
}

func NewIntentUsecase(intentRepo repositories.IntentRepository, pubsubRepo repositories.PubsubChannelRepository) IntentUsecase {
	return &intentUsecase{
		intentRepo: intentRepo,
		pubsubRepo: pubsubRepo,
		log:        logger.New("IntentUsecase"),
	}
}

func (i *intentUsecase) HandleEvents(ctx context.Context, in *entities.IncomingEvent) (err error) {
	log := i.log.WithServiceInfo("HandleEvents").WithRequestID(shio.ReqIDFromContext(ctx)).WithField("channelID", in.ChannelID)
	intent, err := i.intentRepo.Detect(ctx, in)

	if err == repositories.ErrMessageTypeNotSupported {
		intent = &entities.Intent{
			Name:                     in.Message.GetType(),
			Parameters:               in.Message.Parameters,
			FulfillmentText:          "",
			AllRequiredParamsPresent: true,
		}
	} else if err != nil {
		log.WithError(err).Error("detect intent failed")
		return err
	}

	in.Intent = intent
	log.Printf("intent result: %s", intent.Name)

	err = i.pubsubRepo.PublishFulfillmentInput(ctx, in)
	if err != nil {
		log.WithError(err).Error("publish fulfillment failed")
		return err
	}

	// TODO: should we remove this and handle in in fulfillment?
	if intent.FulfillmentText != "" {
		log.Printf("send fulfillment text, %s", intent.FulfillmentText)
		outgoing := &entities.OutgoingEvent{
			RequestID: in.RequestID,
			ChannelID: in.ChannelID,
			Type:      entities.OutgoingEventTypeMessage,
			OutgoingMessage: &entities.OutgoingMessage{
				ReplyToken:  in.ReplyToken,
				Provider:    in.Provider,
				RecipientID: in.Source.GetSourceID(),
				PayloadJson: line2.BuildTextMessage(intent.FulfillmentText),
			}, // TODO assume line for now, to support multi provider
		}

		if e := i.pubsubRepo.PublishOutgoingEventInput(ctx, outgoing); e != nil {
			log.WithError(e).Print("publish outgoing event failed")
		}
	}
	return nil
}
