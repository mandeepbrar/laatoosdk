package testutils

// aguihelpers.go — AG-UI / SSE conversation helpers for integration tests.
//
// ConversationSimulator drives any agent that speaks the AG-UI protocol
// (POST /api/copilotkit + SSE response stream).  It lives in the SDK so
// every solution's integration-test suite can use it without copy-pasting.
//
// Typical usage:
//
//	sim := testutils.NewConversationSimulator(t, sessionID, httpBase, token, authHeader, "")
//	turn1 := sim.SendTurn("Create a solution called foo.")
//	turn2 := sim.SendTurn("Add an education application.")

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// SimTurn holds the output of one AG-UI conversation turn.
type SimTurn struct {
	// Thoughts contains each THOUGHT message text (messageId starts with "thought_").
	Thoughts []string
	// FinalResponse is the concatenated text of all non-thought messages.
	// Most turns produce exactly one FINALRESPONSE message.
	FinalResponse string
	// HITLTask is non-nil when the workflow paused at a manual activity.
	// It is automatically included in the next SendTurn call as state.hitlTask.
	HITLTask map[string]interface{}
	// Error is set when the server emitted a RUN_ERROR event.
	Error string
}

// AllText joins thoughts and final response with newlines.
// Use when you want to judge the full builder output regardless of message type.
func (s SimTurn) AllText() string {
	var parts []string
	parts = append(parts, s.Thoughts...)
	if s.FinalResponse != "" {
		parts = append(parts, s.FinalResponse)
	}
	return strings.Join(parts, "\n")
}

// ConversationSimulator sends turns to an agent via the AG-UI protocol and
// captures the SSE response stream.
//
// Create with NewConversationSimulator; call SendTurn for each user message.
// The HITL task context is automatically carried between turns.
type ConversationSimulator struct {
	t          *testing.T
	httpBase   string // base URL of the engine, e.g. http://localhost:8081
	token      string
	authHeader string
	agentID    string
	sessionID  string
	runNum     int
	prevHITL   map[string]interface{} // round-tripped between turns
}

// NewConversationSimulator creates a simulator that sends turns to agentID on httpBase.
// token and authHeader are used for request authentication (pass "" for unauthenticated).
// agentID is the AG-UI agent name; pass "" to use "artefact_builder".
func NewConversationSimulator(t *testing.T, sessionID, httpBase, token, authHeader, agentID string) *ConversationSimulator {
	t.Helper()
	if authHeader == "" {
		authHeader = "X-Auth-Token"
	}
	if agentID == "" {
		agentID = "artefact_builder"
	}
	return &ConversationSimulator{
		t:          t,
		httpBase:   httpBase,
		token:      token,
		authHeader: authHeader,
		agentID:    agentID,
		sessionID:  sessionID,
	}
}

// SendTurn posts one user message to POST /api/copilotkit and drains the SSE
// response stream until RUN_FINISHED or RUN_ERROR.
//
// The HITL task from the previous turn (if any) is automatically included in
// state.hitlTask so the durable workflow can resume.
func (s *ConversationSimulator) SendTurn(userMessage string) SimTurn {
	s.t.Helper()
	s.runNum++
	runID := fmt.Sprintf("run-%s-%d", s.sessionID, s.runNum)

	innerBody := map[string]interface{}{
		"threadId": s.sessionID,
		"runId":    runID,
		"messages": []map[string]string{
			{"role": "user", "content": userMessage},
		},
	}
	if s.prevHITL != nil {
		innerBody["state"] = map[string]interface{}{"hitlTask": s.prevHITL}
	}

	reqBody := map[string]interface{}{
		"method": "agent/run",
		"params": map[string]string{"agentId": s.agentID},
		"body":   innerBody,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		s.t.Fatalf("ConversationSimulator.SendTurn: marshal: %v", err)
	}

	req, err := http.NewRequest("POST", s.httpBase+"/api/copilotkit", bytes.NewReader(bodyBytes))
	if err != nil {
		s.t.Fatalf("ConversationSimulator.SendTurn: build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	if s.token != "" {
		req.Header.Set(s.authHeader, s.token)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		s.t.Fatalf("ConversationSimulator.SendTurn: POST %s/api/copilotkit: %v", s.httpBase, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.t.Fatalf("ConversationSimulator.SendTurn: HTTP %d: %s", resp.StatusCode, string(body))
	}

	return s.parseSSE(resp.Body)
}

// parseSSE reads AG-UI SSE events until RUN_FINISHED/RUN_ERROR and returns a SimTurn.
//
// AG-UI SSE wire format:  "data: {json}\n\n"
//
// Handled event types:
//
//	TEXT_MESSAGE_START   — begins a message; "thought_" prefix in messageId → THOUGHT
//	TEXT_MESSAGE_CONTENT — appends delta text
//	STATE_SNAPSHOT       — closing snapshot carries hitlTask when workflow pauses
//	RUN_ERROR            — workflow error; stops processing
//	RUN_FINISHED         — normal end of stream
func (s *ConversationSimulator) parseSSE(r io.Reader) SimTurn {
	s.t.Helper()

	msgText := map[string]*strings.Builder{}
	msgIsThought := map[string]bool{}
	var msgOrder []string
	seen := map[string]bool{}

	var result SimTurn

	scanner := bufio.NewScanner(r)
	// Large buffer for config output that can be several KB in one delta
	scanner.Buffer(make([]byte, 64*1024), 512*1024)

	finished := false
	for !finished && scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		var evt map[string]interface{}
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &evt); err != nil {
			s.t.Logf("ConversationSimulator: malformed SSE event (skipped): %v", err)
			continue
		}

		evtType, _ := evt["type"].(string)
		switch evtType {
		case "TEXT_MESSAGE_START":
			msgID, _ := evt["messageId"].(string)
			if msgID != "" && !seen[msgID] {
				seen[msgID] = true
				msgOrder = append(msgOrder, msgID)
				msgIsThought[msgID] = strings.HasPrefix(msgID, "thought_")
				msgText[msgID] = &strings.Builder{}
			}

		case "TEXT_MESSAGE_CONTENT":
			msgID, _ := evt["messageId"].(string)
			delta, _ := evt["delta"].(string)
			if msgID != "" && delta != "" {
				if b, ok := msgText[msgID]; ok {
					b.WriteString(delta)
				} else {
					b = &strings.Builder{}
					b.WriteString(delta)
					msgText[msgID] = b
				}
			}

		case "STATE_SNAPSHOT":
			if snap, ok := evt["snapshot"].(map[string]interface{}); ok {
				if ht, ok := snap["hitlTask"].(map[string]interface{}); ok && len(ht) > 0 {
					result.HITLTask = ht
				}
			}

		case "RUN_ERROR":
			result.Error, _ = evt["message"].(string)
			finished = true

		case "RUN_FINISHED":
			finished = true
		}
	}

	if err := scanner.Err(); err != nil {
		s.t.Logf("ConversationSimulator: SSE scanner error: %v (non-fatal)", err)
	}

	for _, msgID := range msgOrder {
		b, ok := msgText[msgID]
		if !ok {
			continue
		}
		text := b.String()
		if text == "" {
			continue
		}
		if msgIsThought[msgID] {
			result.Thoughts = append(result.Thoughts, text)
		} else {
			if result.FinalResponse == "" {
				result.FinalResponse = text
			} else {
				result.FinalResponse += "\n" + text
			}
		}
	}

	s.prevHITL = result.HITLTask
	return result
}
