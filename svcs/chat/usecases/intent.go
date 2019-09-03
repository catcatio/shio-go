package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/line/line-bot-sdk-go/linebot"
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

	if intent.FulfillmentText != "" {

		msg := linebot.TextMessage{
			Text: intent.FulfillmentText,
		}
		b, _ := json.Marshal([]linebot.Message{&msg})

		e := i.pubsubRepo.PublishOutgoingEventInput(ctx, &entities.OutgoingEvent{
			RequestID: in.RequestID,
			ChannelID: in.ChannelID,
			Type:      entities.OutgoingEventTypeMessage,
			OutgoingMessage: &entities.OutgoingMessage{
				ReplyToken:  in.ReplyToken,
				Provider:    in.Provider,
				RecipientID: in.Source.GetSourceID(),
				PayloadJson: string(b),
			},
		})

		fmt.Println(e)
	}

	return nil
}
