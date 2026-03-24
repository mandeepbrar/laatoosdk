package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// TestServerProcess manages a Laatoo server subprocess started for integration
// testing. Created via StartTestServerProcess.
type TestServerProcess struct {
	// Client is pre-dialled and ready to use after StartTestServerProcess returns.
	Client    *TestServerClient
	cmd       *exec.Cmd
	configDir string
}

// StartTestServerProcess starts the Laatoo server binary as a subprocess with
// the given configDir and waits until the admin RPC port is ready.
//
// serverBinary is the path to the compiled Laatoo server binary. If empty,
// the value of the LAATOO_SERVER_BINARY environment variable is used.
//
// The caller is responsible for calling Stop() when the test suite finishes.
//
// Example (in TestMain):
//
//	proc, err := testutils.StartTestServerProcess("", "test/miniserver")
//	if err != nil { log.Fatal(err) }
//	defer proc.Stop()
//	sharedClient = proc.Client
func StartTestServerProcess(serverBinary, configDir string) (*TestServerProcess, error) {
	if serverBinary == "" {
		serverBinary = os.Getenv("LAATOO_SERVER_BINARY")
	}
	if serverBinary == "" {
		return nil, fmt.Errorf("no server binary: pass a path or set LAATOO_SERVER_BINARY")
	}
	if _, err := os.Stat(serverBinary); err != nil {
		return nil, fmt.Errorf("server binary not found at %q: %w", serverBinary, err)
	}

	cmd := exec.Command(serverBinary, configDir)
	cmd.Env = append(os.Environ(), "LAATOO_TESTMODE=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start server process: %w", err)
	}

	adminAddr := "localhost:20000"
	client, _ := DialTestServer(adminAddr)

	proc := &TestServerProcess{
		Client:    client,
		cmd:       cmd,
		configDir: configDir,
	}

	if err := client.WaitUntilReady(60 * time.Second); err != nil {
		proc.Stop()
		return nil, fmt.Errorf("server did not become ready: %w", err)
	}

	return proc, nil
}

// Stop sends SIGTERM to the server process and waits for it to exit.
func (p *TestServerProcess) Stop() {
	if p.cmd == nil || p.cmd.Process == nil {
		return
	}
	_ = p.cmd.Process.Signal(os.Interrupt)
	_ = p.cmd.Wait()
}
