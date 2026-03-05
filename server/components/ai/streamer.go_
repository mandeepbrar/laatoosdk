package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

const defaultAgentStreamBufferSize = 32

// AgentResponseStreamer streams agent events directly on RequestContext.
type AgentResponseStreamer struct {
	StreamBufferSize int
}

func (streamer *AgentResponseStreamer) SendAgentResponse(ctx core.RequestContext, sessionID string, eventType AgentEventType, data AgentData) error {
	if sessionID == "" {
		sessionID = ctx.CreateUUID()
	}
	if !ctx.IsStreaming() {
		ctx.InitStream(streamer.streamBufferSize())
	}

	payload := utils.StringMap{
		"eventType": string(eventType),
		"sessionId": sessionID,
		"data":      data,
	}

	switch eventType {
	case FINALRESPONSE:
		return ctx.CompleteStream(core.StatusSuccess, payload)
	case ERROR:
		return ctx.CompleteStream(core.StatusInternalError, payload)
	default:
		return ctx.StreamResponse(core.StatusSuccess, payload)
	}
}

func (streamer *AgentResponseStreamer) SendThought(ctx core.RequestContext, sessionID string, data AgentData) error {
	return streamer.SendAgentResponse(ctx, sessionID, THOUGHT, data)
}

func (streamer *AgentResponseStreamer) SendFinalResponse(ctx core.RequestContext, sessionID string, data AgentData) error {
	return streamer.SendAgentResponse(ctx, sessionID, FINALRESPONSE, data)
}

func (streamer *AgentResponseStreamer) SendError(ctx core.RequestContext, sessionID string, data AgentData) error {
	return streamer.SendAgentResponse(ctx, sessionID, ERROR, data)
}

func (streamer *AgentResponseStreamer) streamBufferSize() int {
	if streamer != nil && streamer.StreamBufferSize > 0 {
		return streamer.StreamBufferSize
	}
	return defaultAgentStreamBufferSize
}
