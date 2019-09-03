package entities

//type MessageType string

const (
	MessageTypeTextMessage     string = "textMessage"
	MessageTypeImageMessage    string = "imageMessage"
	MessageTypeVideoMessage    string = "videoMessage"
	MessageTypeAudioMessage    string = "audioMessage"
	MessageTypeFileMessage     string = "fileMessage"
	MessageTypeLocationMessage string = "locationMessage"
	MessageTypeStickerMessage  string = "stickerMessage"
	MessageTypeFollow          string = "follow"
	MessageTypeUnFollow        string = "unfollow"
	MessageTypeJoin            string = "join"
	MessageTypeLeave           string = "leave"
	MessageTypeMemberJoined    string = "memberJoined"
	MessageTypeMemberLeft      string = "memberLeft"
	MessageTypePostback        string = "postback"
	MessageTypeBeacon          string = "beacon"
	MessageTypeAccountLink     string = "accountLink"
	MessageTypeThings          string = "things"
	MessageTypeUnknown         string = "unknown"
)

type EventMessage struct {
	Type       string     `json:"type"`
	Parameters Parameters `json:"parameters"`
}

func (e *EventMessage) GetType() string {
	return e.Type
}

type OutgoingMessage struct {
	ReplyToken  string       `json:"reply_token"`
	Provider    ProviderType `json:"provider"`
	RecipientID string       `json:"recipient_id"`
	PayloadJson string       `json:"payload_json"`
}
