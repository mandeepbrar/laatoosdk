package components

import "laatoo/sdk/server/core"

type Communication struct {
	Subject     string
	Mime        string
	Attachments []string
	Recipients  map[string]string
	Message     []byte
	Info        interface{}
}

type Communicator interface {
	SendCommunication(ctx core.RequestContext, communication *Communication) error
}
