package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/octofoxio/foundation/logger"
)

type PubsubMessageHandlers map[string]PubsubMessageHandler
type PubsubMessageHandler func(ctx context.Context, m pubsub.RawPubsubMessage) error

func NewSendMessagePubsubEndpoint(chat usecases.Chat, line usecases.Line) PubsubMessageHandler {
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

		channelOptions, err := chat.GetChannelConfig(ctx, input.ChannelID)

		if err != nil {
			log.WithError(err).WithField("channelID", input.ChannelID).Error("unable to get channel options")
			return err
		}

		bot, _ := linebot.New(channelOptions.LineChatOptions.ChannelSecret, channelOptions.LineChatOptions.ChannelAccessToken)

		//response, err := bot.ReplyMessage(input.ReplyToken, input.Payload.(*linebot.TextMessage)).Do()
		response, err := bot.ReplyMessage(input.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%v", input.Payload))).Do()

		log.Info(err)
		log.Info(response)

		return chat.SendMessage(ctx, input)
	}
}
