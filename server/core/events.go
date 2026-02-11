package core

import (
	"time"
)

// EventCategory defines the broad classification of an event
type EventCategory string

const (
	BusinessMilestone EventCategory = "business_milestone"
	WorkflowLifecycle EventCategory = "workflow_lifecycle"
	ActivityLifecycle EventCategory = "activity_lifecycle"
	DataAsset         EventCategory = "data_asset"
	SystemAlert       EventCategory = "system_alert"
)

// Event represents a CloudEvents-compliant event structure
// Based on CloudEvents v1.0.2
// Event represents a CloudEvents-compliant event structure
// Based on CloudEvents v1.0.2
type Event struct {
	Message                                // Embeds Data, Tenant, User
	ID              string                 `json:"id"`               // Required
	Source          string                 `json:"source"`           // Required: URI identifying the producer
	SpecVersion     string                 `json:"specversion"`      // Required: "1.0"
	Type            string                 `json:"type"`             // Required: Event type (e.g., com.laatoo.order.placed)
	DataContentType string                 `json:"datacontenttype"`  // Optional: "application/json"
	DataSchema      string                 `json:"dataschema"`       // Optional: URI for the data schema
	Subject         string                 `json:"subject"`          // Optional: Target object (e.g., order-123)
	Time            time.Time              `json:"time"`             // Optional: Event timestamp
	Extensions      map[string]interface{} `json:"extensions"`       // CloudEvents extensions (e.g., tenant, user)
}

// NewEvent creates a basic CloudEvent template
func NewEvent(source, eventType string, data interface{}) *Event {
	return &Event{
		SpecVersion: "1.0",
		Source:      source,
		Type:        eventType,
		Time:        time.Now(),
		Message:     Message{Data: data},
		Extensions:  make(map[string]interface{}),
	}
}
