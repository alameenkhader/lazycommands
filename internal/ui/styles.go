package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Color definitions
	colorGreen  = lipgloss.Color("46")
	colorRed    = lipgloss.Color("196")
	colorYellow = lipgloss.Color("226")
	colorBlue   = lipgloss.Color("39")
	colorGray   = lipgloss.Color("240")
	colorPurple = lipgloss.Color("170")

	// SelectedStyle is used for the currently selected command
	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple)

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	// SuccessStyle is used for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	// RunningStyle is used for currently running commands
	RunningStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	// PendingStyle is used for pending commands
	PendingStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	// SkippedStyle is used for skipped commands
	SkippedStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Italic(true)

	// BorderStyle is used for panel borders
	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorGray)

	// TitleStyle is used for panel titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple)

	// PromptStyle is used for prompts
	PromptStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)
)
