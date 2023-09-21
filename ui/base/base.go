package base

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/styles"
)

type Model struct {
	insert         bool
	urlInput       textinput.Model
	selectedMethod int
	methods        []string
	err            error
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "Enter URL"
	input.CharLimit = 200

	return Model{
		urlInput:       input,
		selectedMethod: 0,
		methods:        []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}
}

func (m Model) GetURL() string {
	return m.urlInput.Value()
}

func (m Model) GetMethod() string {
	return m.methods[m.selectedMethod]
}

func (m *Model) pasteUrl() {
	content, err := clipboard.ReadAll()
	if err == nil {
		m.urlInput.SetValue(content)
	}
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
			m.insert = false
			m.urlInput.Blur()
		case "i":
			m.insert = true
			m.urlInput.Focus()
			return m, nil
		case "p":
			m.pasteUrl()
		case "j", "down":
			if !m.insert {
				m.nextMethod()
			}
		case "k", "up":
			if !m.insert {
				m.prevMethod()
			}
		}
	}

	if m.insert {
		m.urlInput, cmd = m.urlInput.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	methodStyle := lipgloss.NewStyle().Foreground(styles.Orange)
	return lipgloss.JoinHorizontal(lipgloss.Left, "\t"+methodStyle.Render(m.methods[m.selectedMethod])+"\t"+m.urlInput.View())
}
