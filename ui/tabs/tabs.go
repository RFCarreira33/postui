package tabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/config"
	"github.com/rfcarreira33/postui/styles"
)

// Tabs styling
var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	inactiveTab = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(styles.Orange).
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
		focused: config.Base,
		width:   100,
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
	m.focused = (m.focused - 1 + config.NUM_TABS) % config.NUM_TABS
}

func (m *Model) GetFocused() config.Status {
	return m.focused
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
		case "L":
			m.Next()
		case "H":
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
	row := lipgloss.JoinHorizontal(lipgloss.Top, rows...)
	gap := tabGap.Render(strings.Repeat(" ", config.Max(0, m.width-lipgloss.Width(row)-2)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap) + "\n\n"
}
