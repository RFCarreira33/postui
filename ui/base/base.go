package base

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/app/helpers"
	"github.com/rfcarreira33/postui/ui/selector"
)

type Model struct {
	urlInput textinput.Model
	selector selector.Model
	mode     helpers.Mode
	err      error
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "Enter URL"
	input.CharLimit = 200

	return Model{
		urlInput: input,
		selector: *selector.New([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
	}
}

func (m Model) GetURL() string {
	return m.urlInput.Value()
}

func (m *Model) pasteUrl() {
	content, err := clipboard.ReadAll()
	if err == nil {
		m.urlInput.SetValue(content)
	}
}

func (m Model) GetMethod() string {
	return m.selector.GetSelectedItem()
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
			m.mode = helpers.Normal
			m.urlInput.Blur()
		case "i":
			m.mode = helpers.Insert
			m.urlInput.Focus()
			return m, nil
		case "d":
			if m.mode.IsNormal() {
				m.urlInput.SetValue("")
			}
		case "p":
			if !m.mode.IsInsert() {
				m.pasteUrl()
			}
		}
	}

	if m.mode.IsInsert() {
		m.urlInput, cmd = m.urlInput.Update(msg)
	}
	if m.mode.IsNormal() {
		m.selector, cmd = m.selector.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.selector.View(), "\n"+m.urlInput.View())
}
