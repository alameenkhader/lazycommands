package app

import (
	"os"

	"github.com/alameen/lazycommands/internal/executor"
	"github.com/alameen/lazycommands/internal/keys"
	"github.com/alameen/lazycommands/internal/log"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the Bubble Tea application state
type Model struct {
	// Core state
	commands      []*executor.Command
	executing     int               // Index of currently executing command (-1 if none)
	failedCommand *executor.Command // The command that failed (if any)
	workingDir    string            // Current working directory for command execution
	logger        *log.Logger       // Debug logger for command execution

	// UI state
	width   int
	height  int
	ready   bool
	spinner spinner.Model

	// Keyboard
	keys keys.KeyMap
}

// NewModel creates a new Model with the given commands
func NewModel(commands []*executor.Command) Model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = s.Style.Foreground(s.Style.GetForeground())

	// Get initial working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "" // Will use process default
	}

	// Create logger (continue if it fails)
	logger, err := log.NewLogger()
	if err != nil {
		// Log creation failed, but continue without logging
		logger = nil
	}

	return Model{
		commands:      commands,
		executing:     -1,
		failedCommand: nil,
		workingDir:    cwd,
		logger:        logger,
		keys:          keys.DefaultKeyMap(),
		ready:         false,
		spinner:       s,
	}
}

// Init initializes the model and starts command execution
func (m Model) Init() tea.Cmd {
	// Start executing the first command and the spinner
	return tea.Batch(
		(&m).executeNext(),
		m.spinner.Tick,
	)
}

// ExitCode returns the appropriate exit code based on command results
func (m Model) ExitCode() int {
	if m.failedCommand != nil {
		return 1
	}
	for _, cmd := range m.commands {
		if cmd.Status == executor.StatusFailed {
			return 1
		}
	}
	return 0
}

// AllCommandsDone checks if all commands have finished (completed, failed, or skipped)
func (m Model) AllCommandsDone() bool {
	for _, cmd := range m.commands {
		if cmd.Status == executor.StatusPending || cmd.Status == executor.StatusRunning {
			return false
		}
	}
	return true
}

// SkipRemaining marks all pending commands as skipped
func (m *Model) SkipRemaining() {
	for _, cmd := range m.commands {
		if cmd.Status == executor.StatusPending {
			cmd.Status = executor.StatusSkipped

			// Log skipped command
			if m.logger != nil {
				m.logger.LogCommandSkipped(cmd)
			}
		}
	}
}

// Commands returns the list of commands
func (m Model) Commands() []*executor.Command {
	return m.commands
}

// Spinner returns the spinner model
func (m Model) Spinner() spinner.Model {
	return m.spinner
}

// LoggerPath returns the path to the log file, or empty string if no logger
func (m Model) LoggerPath() string {
	if m.logger == nil {
		return ""
	}
	return m.logger.Path()
}

// CloseLogger closes the logger if it exists
func (m *Model) CloseLogger() {
	if m.logger != nil {
		m.logger.Close()
	}
}
