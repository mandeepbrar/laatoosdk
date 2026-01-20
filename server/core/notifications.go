package core

import "laatoo.io/sdk/utils"

type NotificationType string

const (
	INAPP       NotificationType = "INAPP"
	WEBSOCKET   NotificationType = "WEBSOCKET"
	SERVEREVENT NotificationType = "SERVEREVENT"
	AIMESSAGE   NotificationType = "AIMESSAGE"
	EMAIL       NotificationType = "EMAIL"
	SMS         NotificationType = "SMS"
	PUSH        NotificationType = "PUSH"
	WHATSAPP    NotificationType = "WHATSAPP"
	WEBHOOK     NotificationType = "WEBHOOK"
)

type Notification struct {
	NotificationType NotificationType
	Category         string //different types of notifications to be sent on one channel
	Subject          string
	Mime             string
	Attachments      []string
	Recipient        string
	Message          utils.StringMap
	Info             utils.StringMap // Info should contain a string map of Recipients if it has to be sent to multiple users, BroadcastGroup if it has to be sent to some group
}

func ParseNotificationType(str string) (NotificationType, bool) {
	switch str {
	case "INAPP":
		return INAPP, true
	case "SERVEREVENT":
		return SERVEREVENT, true
	case "AIMESSAGE":
		return AIMESSAGE, true
	case "EMAIL":
		return EMAIL, true
	case "SMS":
		return SMS, true
	case "PUSH":
		return PUSH, true
	case "WHATSAPP":
		return WHATSAPP, true
	case "WEBHOOK":
		return WEBHOOK, true
	}
	return "", false
}
