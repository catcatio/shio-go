package entities

import (
	"fmt"
	"time"
)

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

func (p ProviderType) String() string {
	return string(p)
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

func (s *Source) GetSessionID() string {
	return fmt.Sprintf("%s-%s%s", s.UserID, s.RoomID, s.GroupID)
}

type IncomingEvent struct {
	RequestID  string        `json:"request_id"`
	ChannelID  string        `json:"channel_id"`
	Message    *EventMessage `json:"message,omitempty"`
	ReplyToken string        `json:"reply_token"`
	TimeStamp  time.Time     `json:"time_stamp"`
	Source     *Source       `json:"source"`
	Provider   ProviderType  `json:"provider"`
	Original   interface{}   `json:"original"`
	Intent     *Intent       `json:"intent"`
}

type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	PictureUrl  string `json:"picture_url"`
}

type OutgoingEventType string

const (
	OutgoingEventTypeMessage OutgoingEventType = "message"
)

type OutgoingEvent struct {
	RequestID       string            `json:"request_id"`
	ChannelID       string            `json:"channel_id"`
	Type            OutgoingEventType `json:"type"`
	OutgoingMessage *OutgoingMessage  `json:"outgoing_message,omitempty"`
}
