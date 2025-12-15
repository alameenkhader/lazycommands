package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/alameen/lazycommands/internal/executor"
	"github.com/alameen/lazycommands/internal/keys"
)

// Model represents the Bubble Tea application state
type Model struct {
	// Core state
	commands      []*executor.Command
	executing     int // Index of currently executing command (-1 if none)
	failedCommand *executor.Command // The command that failed (if any)

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

	return Model{
		commands:      commands,
		executing:     -1,
		failedCommand: nil,
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
