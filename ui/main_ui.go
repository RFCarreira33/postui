package ui

import (
	"encoding/json"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/config"
	"github.com/rfcarreira33/postui/ui/base"
	"github.com/rfcarreira33/postui/ui/tabs"
	"github.com/rfcarreira33/postui/ui/viewport"
)

type Request struct {
	url     string
	method  string
	headers map[string]string
}

type MainModel struct {
	loaded   bool
	width    int
	height   int
	req      Request
	tabs     tabs.Model
	base     base.Model
	viewport viewport.Model
	insert   bool
	view     bool
	err      error
	quitting bool
}

func New() *MainModel {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	req := Request{
		url:     "",
		method:  "GET",
		headers: headers,
	}

	return &MainModel{
		req:      req,
		tabs:     tabs.New(),
		base:     base.New(),
		viewport: viewport.New(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle key presses
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.loaded = true
			m.width = msg.Width
			m.height = msg.Height - 5
			m.tabs.SetWidth(m.width)
			m.viewport.SetSize(m.width, m.height)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.insert {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			m.insert = false
			m.view = false
		case "i":
			m.insert = true
		case "v":
			m.view = true
		case "R":
			if !m.insert {
				m.req.method = m.base.GetMethod()
				m.req.url = m.base.GetURL()
				curl := exec.Command("curl", "-X", m.req.method, m.req.url)
				for k, v := range m.req.headers {
					curl.Args = append(curl.Args, "-H", k+": "+v)
				}
				output, err := curl.Output()
				if err != nil {
					m.viewport.SetContent("Error running curl check your URL")
					break
				}
				var data interface{}
				json.Unmarshal(output, &data)
				formattedJSON, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					m.viewport.SetContent("Error formatting JSON")
					break
				}
				m.viewport.SetContent(string(formattedJSON))
			}
		}
	}

	var cmd tea.Cmd
	if !m.insert {
		m.tabs, cmd = m.tabs.Update(msg)
	}
	if m.view {
		m.viewport, cmd = m.viewport.Update(msg)
	} else {
		switch m.tabs.GetFocused() {
		default:
			m.base, cmd = m.base.Update(msg)
		}
	}
	return m, cmd
}

func (m MainModel) renderTab() string {
	var tabContent = map[config.Status]string{
		config.Params:  "Params Tab",
		config.Auth:    "Authorization Tab",
		config.Headers: "Headers Tab",
		config.Body:    "Body Tab",
		config.Base:    m.base.View(),
	}

	return lipgloss.PlaceVertical(m.height/3, lipgloss.Top, tabContent[m.tabs.GetFocused()])
}

func (m MainModel) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		sb := strings.Builder{}
		sb.WriteString(m.tabs.View())
		sb.WriteString("\n")
		sb.WriteString(m.renderTab())
		sb.WriteString("\n")
		sb.WriteString("\n")
		sb.WriteString(m.viewport.View())
		return sb.String()
	} else {
		return "loading..."
	}
}
