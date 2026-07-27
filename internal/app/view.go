package app

import (
	"fmt"
	"strings"

	"github.com/alameenkhader/lazycommands/internal/ui"
)

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Always show the command list so statuses stay visible
	b.WriteString(m.renderCommandList())

	// If a command failed, show its output below the list
	if m.failedCommand != nil {
		b.WriteString("\n")
		b.WriteString(m.renderFailedCommandOutput())
	}

	// Add keyboard hints at the bottom
	b.WriteString("\n")
	hint := "q quit"
	if m.LoggerPath() != "" {
		hint += " · o debug log"
	}
	b.WriteString(ui.MutedStyle.Render(hint))

	return b.String()
}

// renderCommandList renders the list of commands with their status icons
func (m Model) renderCommandList() string {
	var b strings.Builder

	for i, cmd := range m.commands {
		// Check if this is the currently running command
		isRunning := (i == m.executing)
		spinnerView := ""
		if isRunning {
			spinnerView = m.spinner.View()
		}
		line := ui.FormatCommandLineWithSpinner(cmd, false, spinnerView)
		b.WriteString(line + "\n")
	}

	return b.String()
}

// renderFailedCommandOutput shows the output of the failed command
func (m Model) renderFailedCommandOutput() string {
	cmd := m.failedCommand

	var b strings.Builder

	b.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed: %s (exit code %d)", cmd.Raw, cmd.ExitCode)) + "\n")

	if cmd.Error != nil {
		b.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("  Error: %s", cmd.Error.Error())) + "\n")
	}

	// Show command output, bounded so the command list above stays visible.
	// Reserve space for the list, the failure header, and the footer hints.
	if len(cmd.Output) > 0 {
		maxLines := m.height - len(m.commands) - 8
		if maxLines < 5 {
			maxLines = 5
		}

		startIdx := 0
		if len(cmd.Output) > maxLines {
			startIdx = len(cmd.Output) - maxLines
		}

		b.WriteString(strings.Repeat("─", 60) + "\n")
		if startIdx > 0 {
			b.WriteString(ui.PendingStyle.Render(fmt.Sprintf("... (%d earlier lines in debug log)", startIdx)) + "\n")
		}
		for i := startIdx; i < len(cmd.Output); i++ {
			b.WriteString(cmd.Output[i] + "\n")
		}
		b.WriteString(strings.Repeat("─", 60) + "\n")
	} else {
		b.WriteString("(No output captured)\n")
	}

	return b.String()
}
