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

func (p *ProviderType) String() string {
	if p == nil {
		return ""
	}

	return string(*p)
}

type Intent struct {
	Name                     string     `json:"name"`
	Parameters               Parameters `json:"parameters"`
	FulfillmentText          string     `json:"fulfillment_text,omitempty"`
	AllRequiredParamsPresent bool       `json:"all_required_params_present"`
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

type event interface {
	getMessage() Message
	getReplyToken() string
	getTimestamp() time.Time
	getSource() *Source
	getProvider() ProviderType
	getOriginalEvent() interface{}
}

type IncomingEvent struct {
	Message     Message
	ReplyToken  string
	TimeStamp   time.Time
	Source      *Source
	Provider    ProviderType
	UserProfile *UserProfile
	Original    interface{}
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
	PictureUrl  string
}
