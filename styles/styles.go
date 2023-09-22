package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	Orange = lipgloss.Color("202")
	Grey   = lipgloss.Color("241")

	ColumnStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Grey)

	FocusedStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Orange)

	TableStyle       = table.DefaultStyles()
	TableHeaderStyle = TableStyle.Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Grey).
				BorderBottom(true).
				Bold(false)

	TableSelectedStyle = TableStyle.Selected.
				Foreground(Orange).
				Bold(false)
)
