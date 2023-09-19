package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	viewport viewport.Model
}

func New() Model {
	return Model{
		viewport: viewport.New(100, 50),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetContent(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) SetSize(w, h int) {
	m.viewport.Width = w - 10
	m.viewport.Height = h/3 + 4
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.viewport.LineDown(1)
		case "k", "up":
			m.viewport.LineUp(1)
		case "ctrl+d":
			m.viewport.ViewDown()
		case "ctrl+u":
			m.viewport.ViewUp()
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}
