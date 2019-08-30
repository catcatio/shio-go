package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
)

type PubsubMessageHandlers map[string]PubsubMessageHandler
type PubsubMessageHandler func(ctx context.Context, m pubsub.RawPubsubMessage) error

func NewSendMessagePubsubEndpoint(sendMessage usecases.SendMessageUsecase) PubsubMessageHandler {
	log := logger.New("SendMessage")
	return func(ctx context.Context, m pubsub.RawPubsubMessage) error {

		if m.Data == nil {
			err := fmt.Errorf("data is nil")
			log.WithError(err).Error("handle pubsub failed")
			return err
		}

		input := new(entities.SendMessageInput)

		if err := json.Unmarshal(m.Data, input); err != nil {
			log.WithError(err).WithField("data", string(m.Data)).Error("unmarshal data failed")
			return err
		}

		log.Println(input)

		return sendMessage.HandleMessage(ctx, input)
	}
}
