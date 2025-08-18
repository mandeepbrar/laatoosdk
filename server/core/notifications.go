package core

import "laatoo.io/sdk/utils"

type NotificationType string

const (
	INAPP    NotificationType = "INAPP"
	EMAIL    NotificationType = "EMAIL"
	SMS      NotificationType = "SMS"
	PUSH     NotificationType = "PUSH"
	WHATSAPP NotificationType = "WHATSAPP"
	WEBHOOK  NotificationType = "WEBHOOK"
)

type Notification struct {
	NotificationType NotificationType
	Subject          string
	Mime             string
	Attachments      []string
	Recipients       map[string]string
	Message          []byte
	Info             utils.StringMap
}

func ParseNotificationType(str string) (NotificationType, bool) {
	switch str {
	case "INAPP":
		return INAPP, true
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
