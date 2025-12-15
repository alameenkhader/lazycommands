package ui

import (
	"github.com/alameen/lazycommands/internal/executor"
)

// StatusIcon returns the icon for a given command status
func StatusIcon(status executor.CommandStatus) string {
	switch status {
	case executor.StatusPending:
		return "⏳"
	case executor.StatusRunning:
		return "▶️ "
	case executor.StatusCompleted:
		return "✔"
	case executor.StatusFailed:
		return "x"
	case executor.StatusSkipped:
		return "⊘ "
	default:
		return "  "
	}
}

// FormatCommandLine formats a command line with its status icon and styling
func FormatCommandLine(cmd *executor.Command, isSelected bool) string {
	icon := StatusIcon(cmd.Status)
	cmdText := cmd.Raw

	// Truncate long commands
	maxLen := 40
	if len(cmdText) > maxLen {
		cmdText = cmdText[:maxLen-3] + "..."
	}

	line := icon + " " + cmdText

	// Apply styling based on status
	switch cmd.Status {
	case executor.StatusRunning:
		line = RunningStyle.Render(line)
	case executor.StatusCompleted:
		line = SuccessStyle.Render(line)
	case executor.StatusFailed:
		line = ErrorStyle.Render(line)
	case executor.StatusSkipped:
		line = SkippedStyle.Render(line)
	case executor.StatusPending:
		line = PendingStyle.Render(line)
	}

	// Highlight if selected
	if isSelected {
		line = SelectedStyle.Render("> " + line)
	} else {
		line = "  " + line
	}

	return line
}

// FormatCommandLineWithSpinner formats a command line with spinner for running commands
func FormatCommandLineWithSpinner(cmd *executor.Command, isSelected bool, spinnerView string) string {
	// Use spinner for running commands, otherwise use status icon
	var icon string
	if cmd.Status == executor.StatusRunning && spinnerView != "" {
		icon = spinnerView
	} else {
		icon = StatusIcon(cmd.Status)
	}

	cmdText := cmd.Raw

	// Truncate long commands
	maxLen := 80
	if len(cmdText) > maxLen {
		cmdText = cmdText[:maxLen-3] + "..."
	}

	line := icon + " " + cmdText

	// Apply styling based on status
	switch cmd.Status {
	case executor.StatusRunning:
		line = RunningStyle.Render(line)
	case executor.StatusCompleted:
		line = SuccessStyle.Render(line)
	case executor.StatusFailed:
		line = ErrorStyle.Render(line)
	case executor.StatusSkipped:
		line = SkippedStyle.Render(line)
	case executor.StatusPending:
		line = PendingStyle.Render(line)
	}

	// Highlight if selected
	if isSelected {
		line = SelectedStyle.Render("> " + line)
	} else {
		line = "  " + line
	}

	return line
}
