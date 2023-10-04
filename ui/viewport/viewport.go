package viewport

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/app/helpers"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
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
	m.viewport.Width = w
	m.viewport.Height = h/3*2 - 3
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
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

func (m Model) headerView() string {
	title := titleStyle.Render("Response")
	line := strings.Repeat("─", helpers.Max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, title)
}

func (m Model) View() string {
	return m.headerView() + "\n" + m.viewport.View()
}
