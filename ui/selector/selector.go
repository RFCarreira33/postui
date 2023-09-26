package selector

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/styles"
)

type Model struct {
	selected int
	items    []string
}

func New(items []string) *Model {
	return &Model{items: items}
}

func (m Model) GetSelectedItem() string {
	return m.items[m.selected]
}

// Move Selected
func (m *Model) next() {
	m.selected = (m.selected + 1) % len(m.items)
}

func (m *Model) prev() {
	length := len(m.items)
	m.selected = (m.selected - 1 + length) % length
}

func (m Model) renderItems() []string {
	selectedStyle := lipgloss.NewStyle().Foreground(styles.Orange)
	var items []string
	for i, item := range m.items {
		if i == m.selected {
			item = selectedStyle.Render(item)
		}
		items = append(items, item+"   ")
	}
	return items
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.prev()
		case "l", "right":
			m.next()
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, m.renderItems()...)
}
