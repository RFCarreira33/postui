package tabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/config"
)

// Tabs styling
var (
	activeTabBorder = lg.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lg.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	inactiveTab = lg.NewStyle().
			Border(tabBorder, true).
			BorderForeground(lg.Color("202")).
			Padding(0, 1)

	activeTab = inactiveTab.Copy().Border(activeTabBorder, true)

	tabGap = inactiveTab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

type Model struct {
	focused config.Status
	width   int
}

func New() Model {
	return Model{
		width: 100,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Move focus
func (m *Model) Next() {
	m.focused = (m.focused + 1) % config.NUM_TABS
}

func (m *Model) Prev() {
	m.focused = (m.focused - 1) % config.NUM_TABS
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l", "right":
			m.Next()
		case "h", "left":
			m.Prev()
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	tabs := []string{"Base", "Params", "Authorization", "Headers", "Body"}
	var rows []string

	for i, tab := range tabs {
		renderedTab := inactiveTab.Render(tab)
		if i == int(m.focused) {
			renderedTab = activeTab.Render(tab)
		}
		rows = append(rows, renderedTab)
	}

	// Join the tabs horizontally
	row := lg.JoinHorizontal(lg.Top, rows...)
	gap := tabGap.Render(strings.Repeat(" ", m.width-62))
	return lg.JoinHorizontal(lg.Bottom, row, gap) + "\n\n"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
