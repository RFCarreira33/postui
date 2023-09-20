package base

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/styles"
)

type Model struct {
	Insert         bool
	UrlInput       textinput.Model
	selectedMethod int
	methods        []string
	err            error
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "Enter URL"
	input.CharLimit = 200

	return Model{
		UrlInput:       input,
		selectedMethod: 0,
		methods:        []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}
}

func (m Model) GetURL() string {
	return m.UrlInput.Value()
}

func (m Model) GetMethod() string {
	return m.methods[m.selectedMethod]
}

func (m *Model) nextMethod() {
	if m.selectedMethod == len(m.methods)-1 {
		return
	}
	m.selectedMethod++
}

func (m *Model) prevMethod() {
	if m.selectedMethod == 0 {
		return
	}
	m.selectedMethod--
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Insert = false
			m.UrlInput.Blur()
		case "i":
			m.Insert = true
			m.UrlInput.Focus()
			return m, nil
		case "j", "down":
			if !m.Insert {
				m.nextMethod()
			}
		case "k", "up":
			if !m.Insert {
				m.prevMethod()
			}
		}
	}

	if m.Insert {
		m.UrlInput, cmd = m.UrlInput.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	methodStyle := lipgloss.NewStyle().Foreground(styles.Orange)
	return lipgloss.JoinHorizontal(lipgloss.Left, "\t"+methodStyle.Render(m.methods[m.selectedMethod])+"\t"+m.UrlInput.View())
}
