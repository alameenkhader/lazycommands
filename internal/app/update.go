package app

import (
	"github.com/alameenkhader/lazycommands/internal/executor"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		// Handle quit
		if key.Matches(msg, m.keys.Quit) {
			// Cancel any running commands
			for _, cmd := range m.commands {
				if cmd.Status == executor.StatusRunning {
					cmd.Cancel()
				}
			}
			return m, tea.Quit
		}

	case executor.TickMsg:
		// Periodic refresh to show streaming output
		// Only keep ticking if a command is running
		if m.executing >= 0 {
			return m, executor.Ticker()
		}
		return m, nil

	case executor.CommandCompletedMsg:
		if msg.Index >= 0 && msg.Index < len(m.commands) {
			cmd := m.commands[msg.Index]

			// Update working directory if changed
			if msg.NewDir != "" {
				m.workingDir = msg.NewDir
			}

			if msg.Error != nil {
				// Command failed - stop execution and show error
				cmd.Status = executor.StatusFailed
				cmd.Error = msg.Error
				cmd.ExitCode = msg.ExitCode
				m.executing = -1
				m.failedCommand = cmd

				// Skip all remaining commands
				(&m).SkipRemaining()

				// Don't quit immediately - let user see the error output
				return m, nil
			}

			// Command succeeded
			cmd.Status = executor.StatusCompleted
			cmd.ExitCode = msg.ExitCode
			m.executing = -1

			// Check if all commands are done
			if m.AllCommandsDone() {
				return m, tea.Quit
			}

			// Execute next command
			return m, (&m).executeNext()
		}
		return m, nil

	default:
		// Handle spinner tick
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
