package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	KeyStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	ValueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	WarnStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
)
