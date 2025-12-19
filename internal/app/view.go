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

	// If there's a failed command, show its output
	if m.failedCommand != nil {
		return m.renderFailedCommandOutput()
	}

	// Otherwise, show the command list
	return m.renderCommandList()
}

// renderCommandList renders the list of commands with their status icons
func (m Model) renderCommandList() string {
	var b strings.Builder

	// b.WriteString(ui.TitleStyle.Render("LazyCommands") + "\n\n")

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

	// Add keyboard hints at the bottom
	b.WriteString("\n")
	b.WriteString(ui.PendingStyle.Render("Press q to quit"))

	// Show log file path if available
	if logPath := m.LoggerPath(); logPath != "" {
		b.WriteString("\n")
		b.WriteString(ui.PendingStyle.Render(fmt.Sprintf("Debug log: %s", logPath)))
	}

	return b.String()
}

// renderFailedCommandOutput shows the output of the failed command
func (m Model) renderFailedCommandOutput() string {
	cmd := m.failedCommand

	var b strings.Builder

	b.WriteString(ui.ErrorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n")
	b.WriteString(ui.ErrorStyle.Render("Command Failed!") + "\n")
	b.WriteString(ui.ErrorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n\n")

	b.WriteString(fmt.Sprintf("Command: %s\n", ui.ErrorStyle.Render(cmd.Raw)))
	b.WriteString(fmt.Sprintf("Exit Code: %s\n\n", ui.ErrorStyle.Render(fmt.Sprintf("%d", cmd.ExitCode))))

	if cmd.Error != nil {
		b.WriteString(fmt.Sprintf("Error: %s\n\n", ui.ErrorStyle.Render(cmd.Error.Error())))
	}

	// Show command output
	if len(cmd.Output) > 0 {
		b.WriteString(ui.TitleStyle.Render("Output:") + "\n")
		b.WriteString(strings.Repeat("─", 60) + "\n")

		// Show all output (or last N lines if too long)
		maxLines := m.height - 15 // Reserve space for header
		if maxLines < 20 {
			maxLines = 20
		}

		startIdx := 0
		if len(cmd.Output) > maxLines {
			startIdx = len(cmd.Output) - maxLines
		}

		for i := startIdx; i < len(cmd.Output); i++ {
			b.WriteString(cmd.Output[i] + "\n")
		}

		// Show indicator if there's more output
		if startIdx > 0 {
			b.WriteString(ui.PendingStyle.Render(fmt.Sprintf("\n... (%d more lines above)", startIdx)))
		}
	} else {
		b.WriteString("(No output captured)\n")
	}

	b.WriteString("\n" + ui.ErrorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n")

	// Show log file path if available
	if logPath := m.LoggerPath(); logPath != "" {
		b.WriteString("\n")
		b.WriteString(ui.PendingStyle.Render(fmt.Sprintf("Full debug log: %s", logPath)))
	}

	return b.String()
}
