package entities

import (
	"encoding/base64"
	"github.com/line/line-bot-sdk-go/linebot"
	"time"
)

type LineEvent struct {
	*linebot.Event
}

func (l *LineEvent) getMessage() *EventMessage {
	return makeMessage(l.Event)
}

func (l *LineEvent) getReplyToken() string {
	return l.ReplyToken
}

func (l *LineEvent) getTimestamp() time.Time {
	return l.Timestamp
}

func (l *LineEvent) getSource() *Source {
	switch l.Source.Type {
	case linebot.EventSourceTypeGroup:
		return &Source{
			Type:    SourceTypeGroup,
			UserID:  l.Source.UserID,
			GroupID: l.Source.GroupID,
		}
	case linebot.EventSourceTypeRoom:
		return &Source{
			Type:   SourceTypeRoom,
			UserID: l.Source.UserID,
			RoomID: l.Source.RoomID,
		}
	}

	// return default as type user
	return &Source{
		Type:   SourceTypeUser,
		UserID: l.Source.UserID,
	}
}

func (l *LineEvent) getProvider() ProviderType {
	return ProviderTypeLine
}

func (l *LineEvent) getOriginalEvent() interface{} {
	return *l
}

func makeMessage(event *linebot.Event) *EventMessage {
	switch event.Type {
	case linebot.EventTypeFollow:
		return &EventMessage{
			Type:       MessageTypeFollow,
			Parameters: Parameters{}, // empty
		}
	case linebot.EventTypeUnfollow:
		return &EventMessage{
			Type:       MessageTypeUnFollow,
			Parameters: Parameters{}, // empty
		}
	case linebot.EventTypeJoin:
		return &EventMessage{
			Type:       MessageTypeJoin,
			Parameters: Parameters{}, // empty
		}
	case linebot.EventTypeLeave:
		return &EventMessage{
			Type:       MessageTypeLeave,
			Parameters: Parameters{}, // empty
		}
	case linebot.EventTypeMemberJoined:
		return &EventMessage{
			Type: MessageTypeMemberJoined,
			Parameters: Parameters{
				"Members": event.Joined,
			},
		}
	case linebot.EventTypeMemberLeft:
		return &EventMessage{
			Type: MessageTypeMemberLeft,
			Parameters: Parameters{
				"Members": event.Left,
			},
		}
	case linebot.EventTypePostback:
		return &EventMessage{
			Type: MessageTypePostback,
			Parameters: Parameters{
				"Data":   event.Postback.Data,
				"Params": event.Postback.Params,
			},
		}
	case linebot.EventTypeBeacon:
		return &EventMessage{
			Type: MessageTypeBeacon,
			Parameters: Parameters{
				"Hwid":          event.Beacon.Hwid,
				"DeviceMessage": base64.StdEncoding.EncodeToString(event.Beacon.DeviceMessage),
			},
		}
	case linebot.EventTypeAccountLink:
		return &EventMessage{
			Type: MessageTypeAccountLink,
			Parameters: Parameters{
				"Result": event.AccountLink.Result,
				"Nonce":  event.AccountLink.Nonce,
			},
		}
	case linebot.EventTypeThings:
		return &EventMessage{
			Type: MessageTypeThings,
			Parameters: Parameters{
				"DeviceID": event.Things.DeviceID,
				"Type":     event.Things.Type,
			},
		}
	}

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		return &EventMessage{
			Type: MessageTypeTextMessage,
			Parameters: Parameters{
				"ID":   message.ID,
				"Text": message.Text,
			},
		}
	case *linebot.ImageMessage:
		return &EventMessage{
			Type: MessageTypeImageMessage,
			Parameters: Parameters{
				"ID":                 message.ID,
				"OriginalContentURL": message.OriginalContentURL,
				"PreviewImageURL":    message.PreviewImageURL,
			},
		}
	case *linebot.VideoMessage:
		return &EventMessage{
			Type: MessageTypeVideoMessage,
			Parameters: Parameters{
				"ID":                 message.ID,
				"OriginalContentURL": message.OriginalContentURL,
				"PreviewImageURL":    message.PreviewImageURL,
			},
		}
	case *linebot.AudioMessage:
		return &EventMessage{
			Type: MessageTypeAudioMessage,
			Parameters: Parameters{
				"ID":                 message.ID,
				"OriginalContentURL": message.OriginalContentURL,
				"Duration":           message.Duration,
			},
		}
	case *linebot.FileMessage:
		return &EventMessage{
			Type: MessageTypeFileMessage,
			Parameters: Parameters{
				"ID":       message.ID,
				"FileName": message.FileName,
				"FileSize": message.FileSize,
			},
		}
	case *linebot.LocationMessage:
		return &EventMessage{
			Type: MessageTypeLocationMessage,
			Parameters: Parameters{
				"ID":        message.ID,
				"Title":     message.Title,
				"Address":   message.Address,
				"Latitude":  message.Latitude,
				"Longitude": message.Longitude,
			},
		}
	case *linebot.StickerMessage:
		return &EventMessage{
			Type: MessageTypeStickerMessage,
			Parameters: Parameters{
				"ID":        message.ID,
				"PackageID": message.PackageID,
				"StickerID": message.StickerID,
			},
		}
	default:
		return &EventMessage{
			Type:       MessageTypeUnknown,
			Parameters: Parameters{},
		}
	}
}

func (l *LineEvent) IncomingEvent(requestID, channelID string) *IncomingEvent {
	return &IncomingEvent{
		requestID,
		channelID,
		l.getMessage(),
		l.getReplyToken(),
		l.getTimestamp(),
		l.getSource(),
		l.getProvider(),
		l.getOriginalEvent(),
		nil,
	}
}
