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
	UserID  string     `json:"user_id,omitempty"`
	GroupID string     `json:"group_id,omitempty"`
	RoomID  string     `json:"room_id,omitempty"`
}

func (s *Source) GetUserID() string {
	return s.UserID
}

func (s *Source) GetSourceID() string {
	if s.GroupID != "" {
		return s.GroupID
	}

	if s.RoomID != "" {
		return s.RoomID
	}

	return s.UserID
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
	RequestID   string       `json:"request_id"`
	Message     Message      `json:"message"`
	ReplyToken  string       `json:"reply_token"`
	TimeStamp   time.Time    `json:"timestamp"`
	Source      *Source      `json:"source,omitempty"`
	Provider    ProviderType `json:"provider"`
	UserProfile *UserProfile `json:"user_profile,omitempty"`
	Original    interface{}  `json:"original"`
	Intent      *Intent      `json:"intent"`
}

type UserProfile struct {
	ID          string
	DisplayName string
	PictureUrl  string
}

type SendMessageInput struct {
	RequestID   string       `json:"request_id"`
	ChannelID   string       `json:"channel_id"`
	ReplyToken  string       `json:"reply_token"`
	Provider    ProviderType `json:"provider"`
	RecipientID string       `json:"recipient_id"`
	Payload     interface{}  `json:"payload"`
}
