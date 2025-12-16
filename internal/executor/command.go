package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CommandStatus represents the execution state of a command
type CommandStatus int

const (
	StatusPending CommandStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusSkipped
)

// String returns a string representation of the command status
func (s CommandStatus) String() string {
	switch s {
	case StatusPending:
		return "Pending"
	case StatusRunning:
		return "Running"
	case StatusCompleted:
		return "Completed"
	case StatusFailed:
		return "Failed"
	case StatusSkipped:
		return "Skipped"
	default:
		return "Unknown"
	}
}

// Command wraps a shell command with its execution state
type Command struct {
	ID          int
	Raw         string        // Original command string
	Status      CommandStatus // Current execution status
	Output      []string      // Captured stdout/stderr lines
	ExitCode    int           // Exit code of the command
	StartTime   time.Time     // When the command started
	EndTime     time.Time     // When the command finished
	Error       error         // Error if the command failed
	WorkingDir  string        // Working directory for this command
	IsCdCommand bool          // True if this is a cd command
	ctx         context.Context
	cancel      context.CancelFunc
}

const maxOutputLines = 1000

// NewCommand creates a new Command instance
func NewCommand(id int, raw string) *Command {
	ctx, cancel := context.WithCancel(context.Background())
	return &Command{
		ID:     id,
		Raw:    raw,
		Status: StatusPending,
		Output: make([]string, 0, maxOutputLines),
		ctx:    ctx,
		cancel: cancel,
	}
}

// AppendOutput adds a line to the command's output, maintaining a sliding window
// to prevent memory issues with very long outputs
func (c *Command) AppendOutput(line string) {
	c.Output = append(c.Output, line)
	if len(c.Output) > maxOutputLines {
		// Keep only the last maxOutputLines
		c.Output = c.Output[len(c.Output)-maxOutputLines:]
	}
}

// Cancel cancels the command's context
func (c *Command) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Duration returns the duration of the command execution
func (c *Command) Duration() time.Duration {
	if c.StartTime.IsZero() {
		return 0
	}
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// ParseCdCommand checks if a command is a cd command and extracts the target directory.
// Only supports absolute paths (no relative paths, no ~).
// Returns: (isCd, targetDir, error)
func ParseCdCommand(cmdStr string) (bool, string, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(cmdStr)

	// Check if it starts with "cd "
	if !strings.HasPrefix(trimmed, "cd ") {
		return false, "", nil
	}

	// Extract the directory argument
	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		return true, "", fmt.Errorf("cd: missing directory argument")
	}
	if len(parts) > 2 {
		return true, "", fmt.Errorf("cd: too many arguments")
	}

	targetDir := parts[1]

	// Remove quotes if present
	targetDir = strings.Trim(targetDir, `"'`)

	// Validate it's an absolute path
	if !filepath.IsAbs(targetDir) {
		return true, "", fmt.Errorf("cd: only absolute paths are supported (got: %s)", targetDir)
	}

	// Check if directory exists
	info, err := os.Stat(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return true, "", fmt.Errorf("cd: directory does not exist: %s", targetDir)
		}
		return true, "", fmt.Errorf("cd: cannot access directory: %v", err)
	}

	if !info.IsDir() {
		return true, "", fmt.Errorf("cd: not a directory: %s", targetDir)
	}

	return true, targetDir, nil
}
