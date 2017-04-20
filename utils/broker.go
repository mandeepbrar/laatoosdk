package utils

import (
	"github.com/pmylund/go-cache"
)

type EventHandler func(map[string]interface{}) error

var (
	registeredHandlers = cache.New(cache.NoExpiration, 0)
)

type Event struct {
	Type      string
	EventData map[string]interface{}
}

func RegisterEventHandler(eventType string, handler EventHandler) {
	eventHandlersInt, prs := registeredHandlers.Get(eventType)
	var eventHandlers []EventHandler
	if !prs {
		eventHandlers = []EventHandler{}
	} else {
		eventHandlers = eventHandlersInt.([]EventHandler)
	}
	eventHandlers = append(eventHandlers, handler)
	registeredHandlers.Set(eventType, eventHandlers, cache.NoExpiration)
}

func FireEvent(event *Event) error {
	handlersInt, present := registeredHandlers.Get(event.Type)
	if present {
		handlers := handlersInt.([]EventHandler)
		for _, handler := range handlers {
			err := handler(event.EventData)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
