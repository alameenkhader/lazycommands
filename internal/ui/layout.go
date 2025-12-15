package ui

import "github.com/charmbracelet/lipgloss"

// Layout manages the split-panel layout
type Layout struct {
	Width  int
	Height int
}

// NewLayout creates a new Layout with the given dimensions
func NewLayout(width, height int) Layout {
	return Layout{
		Width:  width,
		Height: height,
	}
}

// LeftWidth returns the width of the left panel (33% of total width)
func (l Layout) LeftWidth() int {
	leftWidth := l.Width / 3
	if leftWidth < 30 {
		leftWidth = 30
	}
	return leftWidth
}

// RightWidth returns the width of the right panel
func (l Layout) RightWidth() int {
	rightWidth := l.Width - l.LeftWidth() - 3 // -3 for border and spacing
	if rightWidth < 40 {
		rightWidth = 40
	}
	return rightWidth
}

// PanelHeight returns the height available for content (minus borders)
func (l Layout) PanelHeight() int {
	return l.Height - 2 // -2 for top and bottom borders/spacing
}

// Render combines left and right panel content into a single view
func (l Layout) Render(left, right string) string {
	leftPanel := lipgloss.NewStyle().
		Width(l.LeftWidth()).
		Height(l.Height).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("240")).
		Render(left)

	rightPanel := lipgloss.NewStyle().
		Width(l.RightWidth()).
		Height(l.Height).
		Padding(0, 1).
		Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}
