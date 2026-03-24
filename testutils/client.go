// Package testutils provides helpers for writing integration tests against a
// running Laatoo server. It is importable by any plugin (they all depend on the
// SDK) without creating a circular dependency with the server package.
//
// Usage pattern:
//
//	func TestMain(m *testing.M) {
//	    addr := os.Getenv("LAATOO_TEST_SERVER") // e.g. "localhost:20000"
//	    if addr == "" {
//	        fmt.Println("skipping: set LAATOO_TEST_SERVER=<adminAddr>")
//	        os.Exit(0)
//	    }
//	    client, err := testutils.DialTestServer(addr)
//	    if err != nil { log.Fatal(err) }
//	    if err := client.WaitUntilReady(30 * time.Second); err != nil { log.Fatal(err) }
//	    // share client with tests via package-level var
//	    os.Exit(m.Run())
//	}
package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestServerClient communicates with the Laatoo admin RPC port (:20000) using
// the JSON-RPC-over-HTTP protocol that standalone.go exposes.
type TestServerClient struct {
	adminAddr string // e.g. "localhost:20000"
	httpBase  string // e.g. "http://localhost:8080" (set via SetHTTPBase)
	reqID     int
}

// DialTestServer returns a client connected to the admin RPC port at adminAddr.
// adminAddr should be "host:port", e.g. "localhost:20000".
// The server does not need to be ready yet; call WaitUntilReady before tests.
func DialTestServer(adminAddr string) (*TestServerClient, error) {
	return &TestServerClient{adminAddr: adminAddr}, nil
}

// SetHTTPBase configures the base URL for HTTP calls made via InvokeHTTP.
func (c *TestServerClient) SetHTTPBase(base string) { c.httpBase = base }

// WaitUntilReady polls Ping until the server responds or the timeout elapses.
func (c *TestServerClient) WaitUntilReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		var reply struct{ Status string }
		lastErr = c.call("TestOps.Ping", struct{}{}, &reply)
		if lastErr == nil && reply.Status == "ok" {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("server not ready after %v: %w", timeout, lastErr)
}

// ── Activity invocation ───────────────────────────────────────────────────────

type InvokeActivityResult struct {
	Result   map[string]string
	Messages []string
	Error    string
}

// InvokeActivity runs a named activity on the server with the given params.
// Streamed messages (THOUGHT / FINALRESPONSE) are returned in Messages.
// If sessionID is non-empty they are also stored server-side for later
// retrieval via GetStreamMessages.
func (c *TestServerClient) InvokeActivity(name string, params map[string]string, sessionID string) (*InvokeActivityResult, error) {
	args := map[string]interface{}{
		"ActivityName": name,
		"Params":       params,
		"SessionID":    sessionID,
	}
	var reply InvokeActivityResult
	if err := c.call("TestOps.InvokeActivity", args, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

// ── HITL helpers ──────────────────────────────────────────────────────────────

// HITLTaskRef is a reference to a pending HITL task on the server.
type HITLTaskRef struct {
	TaskID     string
	WorkflowID string
	InstanceID string
	ActivityID string
}

// GetPendingHITLTasks returns all pending HITL tasks for a session.
func (c *TestServerClient) GetPendingHITLTasks(sessionID string) ([]HITLTaskRef, error) {
	args := map[string]string{"SessionID": sessionID}
	var reply struct{ Tasks []HITLTaskRef }
	if err := c.call("TestOps.GetPendingHITLTasks", args, &reply); err != nil {
		return nil, err
	}
	return reply.Tasks, nil
}

// WaitForHITL polls until at least one HITL task appears for the session.
// timeoutSec defaults to 30 if <= 0.
func (c *TestServerClient) WaitForHITL(sessionID string, timeoutSec int) ([]HITLTaskRef, error) {
	args := map[string]interface{}{"SessionID": sessionID, "TimeoutSec": timeoutSec}
	var reply struct {
		Tasks []HITLTaskRef
		Error string
	}
	if err := c.call("TestOps.WaitForHITL", args, &reply); err != nil {
		return nil, err
	}
	if reply.Error != "" {
		return nil, fmt.Errorf("WaitForHITL: %s", reply.Error)
	}
	return reply.Tasks, nil
}

// SignalHITL delivers message to the HITL task identified by taskID.
func (c *TestServerClient) SignalHITL(taskID, sessionID, message string) error {
	args := map[string]string{
		"TaskID":    taskID,
		"SessionID": sessionID,
		"Message":   message,
	}
	var reply struct{ Error string }
	if err := c.call("TestOps.SignalHITL", args, &reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return fmt.Errorf("SignalHITL: %s", reply.Error)
	}
	return nil
}

// ── Stream message capture ────────────────────────────────────────────────────

// GetStreamMessages drains and returns all captured stream messages for the session.
func (c *TestServerClient) GetStreamMessages(sessionID string) ([]string, error) {
	args := map[string]string{"SessionID": sessionID}
	var reply struct{ Messages []string }
	if err := c.call("TestOps.GetStreamMessages", args, &reply); err != nil {
		return nil, err
	}
	return reply.Messages, nil
}

// ── HTTP helpers ──────────────────────────────────────────────────────────────

// InvokeHTTP makes an HTTP request to the server's HTTP engine.
// Requires SetHTTPBase to have been called first.
func (c *TestServerClient) InvokeHTTP(method, path, contentType, body string) (*http.Response, error) {
	if c.httpBase == "" {
		return nil, fmt.Errorf("httpBase not set: call SetHTTPBase first")
	}
	url := c.httpBase + path
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return http.DefaultClient.Do(req)
}

// ── JSON-RPC transport ────────────────────────────────────────────────────────

type rpcRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	ID     int         `json:"id"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  interface{}     `json:"error"`
	ID     int             `json:"id"`
}

func (c *TestServerClient) call(method string, args interface{}, result interface{}) error {
	c.reqID++
	reqBody, err := json.Marshal(rpcRequest{Method: method, Params: []interface{}{args}, ID: c.reqID})
	if err != nil {
		return fmt.Errorf("marshal RPC request: %w", err)
	}

	url := "http://" + c.adminAddr + "/"
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("admin RPC call %s: %w", method, err)
	}
	defer resp.Body.Close()

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("decode RPC response: %w", err)
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("RPC error from %s: %v", method, rpcResp.Error)
	}
	if result != nil && rpcResp.Result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("unmarshal RPC result: %w", err)
		}
	}
	return nil
}
