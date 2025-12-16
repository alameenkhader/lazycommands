package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alameen/lazycommands/internal/executor"
)

// Logger handles writing command execution logs to a temporary file
type Logger struct {
	file *os.File
	mu   sync.Mutex
	path string
}

// NewLogger creates a new logger that writes to a file in the system temp directory
func NewLogger() (*Logger, error) {
	// Generate log file name with timestamp and PID
	timestamp := time.Now().Format("2006-01-02-150405")
	pid := os.Getpid()
	filename := fmt.Sprintf("lazycommands-%s-%d.log", timestamp, pid)

	// Create log file in temp directory
	logPath := filepath.Join(os.TempDir(), filename)
	file, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	logger := &Logger{
		file: file,
		path: logPath,
	}

	// Write header
	logger.writeHeader()

	return logger, nil
}

// writeHeader writes the log file header
func (l *Logger) writeHeader() {
	l.mu.Lock()
	defer l.mu.Unlock()

	header := fmt.Sprintf("LazyCommands Execution Log\nStarted: %s\n%s\n\n",
		time.Now().Format("2006-01-02 15:04:05"),
		"=======================================================")
	l.file.WriteString(header)
	l.file.Sync()
}

// LogCommandStart logs the start of a command execution
func (l *Logger) LogCommandStart(cmd *executor.Command) {
	if l == nil || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	workingDir := cmd.WorkingDir
	if workingDir == "" {
		workingDir = "(default)"
	}

	entry := fmt.Sprintf("[%s] [CMD-%d] START: %s (WorkingDir: %s)\n",
		timestamp, cmd.ID, cmd.Raw, workingDir)
	l.file.WriteString(entry)
	l.file.Sync()
}

// LogCommandOutput logs a line of command output
func (l *Logger) LogCommandOutput(cmdID int, line string) {
	if l == nil || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	entry := fmt.Sprintf("[%s] [CMD-%d] OUTPUT: %s\n", timestamp, cmdID, line)
	l.file.WriteString(entry)
	l.file.Sync()
}

// LogCommandEnd logs the completion of a command
func (l *Logger) LogCommandEnd(cmd *executor.Command) {
	if l == nil || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	duration := cmd.Duration()

	entry := fmt.Sprintf("[%s] [CMD-%d] END: exit_code=%d duration=%v status=%s",
		timestamp, cmd.ID, cmd.ExitCode, duration, cmd.Status)

	if cmd.Error != nil {
		entry += fmt.Sprintf(" error=\"%v\"", cmd.Error)
	}

	entry += "\n\n"
	l.file.WriteString(entry)
	l.file.Sync()
}

// LogCommandSkipped logs when a command is skipped
func (l *Logger) LogCommandSkipped(cmd *executor.Command) {
	if l == nil || l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	entry := fmt.Sprintf("[%s] [CMD-%d] SKIPPED: %s\n", timestamp, cmd.ID, cmd.Raw)
	l.file.WriteString(entry)
	l.file.Sync()
}

// Path returns the path to the log file
func (l *Logger) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

// Close closes the log file
func (l *Logger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Write footer
	footer := fmt.Sprintf("\n%s\nCompleted: %s\n",
		"=======================================================",
		time.Now().Format("2006-01-02 15:04:05"))
	l.file.WriteString(footer)
	l.file.Sync()

	return l.file.Close()
}
