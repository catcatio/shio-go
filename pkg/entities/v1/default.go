package entities

import "time"

type Parameters map[string]interface{}

func (p Parameters) String(name string) string {
	if _, ok := p[name]; !ok {
		return ""
	}

	if value, ok := p[name].(string); ok {
		return value
	}

	if value, ok := p[name].(*string); ok {
		return *value
	}

	panic("invalid data type")
}

type SourceType string

const (
	SourceTypeUser  SourceType = "user"
	SourceTypeGroup SourceType = "group"
	SourceTypeRoom  SourceType = "room"
)

type ProviderType string

const (
	ProviderTypeLine    ProviderType = "line"
	ProviderTypeUnknown ProviderType = "unknown"
)

type Intent struct {
	Name       string     `json:"name"`
	Parameters Parameters `json:"parameters"`
}

type Fulfillment struct {
	Name       string     `json:"name"`
	Parameters Parameters `json:"parameters"`
}

type Source struct {
	Type    SourceType `json:"type"`
	UserID  string     `json:"userId,omitempty"`
	GroupID string     `json:"groupId,omitempty"`
	RoomID  string     `json:"roomId,omitempty"`
}

type OutgoingEvent struct {
}

type Event interface {
	GetMessage() Message
	GetReplyToken() string
	GetTimestamp() time.Time
	GetSource() *Source
	GetProvider() ProviderType
	GetOriginalEvent() interface{}
}

type ParsedEvent struct {
	RequestID    string       `json:"requestId"`
	Message      Message      `json:"message"`
	ReplyToken   string       `json:"replyToken"`
	TimeStamp    time.Time    `json:"timestamp"`
	Source       *Source      `json:"source,omitempty"`
	ProviderType ProviderType `json:"providerType"`
	UserProfile  *UserProfile `json:"userProfile,omitempty"`
	Original     interface{}  `json:"original"`
	Intent       *Intent      `json:"intent"`
}

type UserProfile struct {
	ID          string
	DisplayName string
	PictureUrl  *string
}
