package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/alameen/lazycommands/internal/executor"
)

// executeNext starts executing the next pending command
func (m *Model) executeNext() tea.Cmd {
	// Find the next pending command
	for i, cmd := range m.commands {
		if cmd.Status == executor.StatusPending {
			// Update executing index
			m.executing = i
			// Return batch: start execution + start ticker for UI refresh
			return tea.Batch(
				executor.ExecuteCommand(i, cmd),
				executor.Ticker(),
			)
		}
	}

	// No more commands to execute
	return nil
}
