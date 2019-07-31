package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"net/http"
)

type requestParser struct {
	channelSecret string
}

func (p *requestParser) Parse(r *http.Request) (events []*linebot.Event, destination string, err error) {
	return parse(p.channelSecret, r)
}

type RequestParser interface {
	Parse(r *http.Request) (events []*linebot.Event, destination string, err error)
}

func NewParser(channelSecret string) RequestParser {
	return &requestParser{
		channelSecret: channelSecret,
	}
}

func parse(channelSecret string, r *http.Request) (events []*linebot.Event, destination string, err error) {
	defer func() {
		_ = r.Body.Close()
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	if len(r.Header.Get("x-shio-debug")) <= 0 && !validateSignature(channelSecret, r.Header.Get("X-Line-Signature"), body) {
		err = linebot.ErrInvalidSignature
		return
	}

	request := &struct {
		Events      []*linebot.Event `json:"events"`
		Destination string           `json:"destination"`
	}{}

	if err = json.Unmarshal(body, request); err != nil {
		return
	}
	return request.Events, request.Destination, nil
}

func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}
