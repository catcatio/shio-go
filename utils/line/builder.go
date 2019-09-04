package line

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
)

func BuildTextMessage(texts ...string) string {
	msgs := make([]*linebot.TextMessage, 0)
	for _, t := range texts {
		msgs = append(msgs, linebot.NewTextMessage(t))
	}

	b, _ := json.Marshal(msgs)

	return string(b)
}
