package dev

import (
	"os/exec"
	"testing"
)

// captureCmd records the arguments passed to execCommand without actually
// running a subprocess.
type captureCmd struct {
	path string
	args []string
}

// TestExecuteWithClaudeNoShell verifies that executeWithClaude does NOT use
// "bash" as the subprocess — i.e. it must invoke the claude binary directly
// rather than writing a shell script and passing it through bash.
//
// RED state: current executeWithClaude calls exec.Command("bash", scriptFile).
// This test MUST fail until FOUND-02 is fixed.
func TestExecuteWithClaudeNoShell(t *testing.T) {
	var captured captureCmd
	ran := false

	// Override execCommand so no real subprocess is started.
	orig := execCommand
	execCommand = func(name string, arg ...string) *exec.Cmd {
		captured.path = name
		captured.args = append([]string{name}, arg...)
		ran = true
		// Return a no-op cmd that will succeed immediately.
		return exec.Command("true")
	}
	defer func() { execCommand = orig }()

	// Set a dummy API key to bypass the early-return guard.
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	err := executeWithClaude("test prompt", "test-task")
	if err != nil {
		// If the injected "true" binary isn't available, tolerate errors, but
		// still check the captured command name.
		t.Logf("executeWithClaude returned error (may be acceptable in CI): %v", err)
	}

	if !ran {
		t.Fatal("execCommand was never called; executeWithClaude must invoke execCommand to run a subprocess")
	}

	if captured.path == "bash" {
		t.Errorf("executeWithClaude used 'bash' as the command; want direct invocation of 'claude' (FOUND-02: shell injection via bash wrapper)")
	}

	// Assert no argument is a shell script path.
	for _, arg := range captured.args[1:] {
		if len(arg) > 0 && arg[0] == '/' {
			t.Errorf("executeWithClaude passed a file path argument %q to the subprocess; want prompt passed via stdin instead (FOUND-02)", arg)
		}
	}
}

// TestExecuteWithOpenCodeNoShell verifies that executeWithOpenCode does NOT
// use "bash" as the subprocess.
//
// RED state: current executeWithOpenCode calls exec.Command("bash", scriptFile).
// This test MUST fail until FOUND-02 is fixed.
func TestExecuteWithOpenCodeNoShell(t *testing.T) {
	var captured captureCmd
	ran := false

	orig := execCommand
	execCommand = func(name string, arg ...string) *exec.Cmd {
		captured.path = name
		captured.args = append([]string{name}, arg...)
		ran = true
		return exec.Command("true")
	}
	defer func() { execCommand = orig }()

	err := executeWithOpenCode("test prompt", "test-task")
	if err != nil {
		t.Logf("executeWithOpenCode returned error (may be acceptable in CI): %v", err)
	}

	if !ran {
		t.Fatal("execCommand was never called; executeWithOpenCode must invoke execCommand to run a subprocess")
	}

	if captured.path == "bash" {
		t.Errorf("executeWithOpenCode used 'bash' as the command; want direct invocation of 'opencode' (FOUND-02: shell injection via bash wrapper)")
	}

	for _, arg := range captured.args[1:] {
		if len(arg) > 0 && arg[0] == '/' {
			t.Errorf("executeWithOpenCode passed a file path argument %q to the subprocess; want prompt via stdin (FOUND-02)", arg)
		}
	}
}
