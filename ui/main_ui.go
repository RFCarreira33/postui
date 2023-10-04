package ui

import (
	"encoding/json"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/app/helpers"
	"github.com/rfcarreira33/postui/app/request"
	"github.com/rfcarreira33/postui/styles"
	"github.com/rfcarreira33/postui/ui/base"
	"github.com/rfcarreira33/postui/ui/body"
	"github.com/rfcarreira33/postui/ui/headers"
	"github.com/rfcarreira33/postui/ui/params"
	"github.com/rfcarreira33/postui/ui/tabs"
	"github.com/rfcarreira33/postui/ui/viewport"
)

type MainModel struct {
	loaded   bool
	width    int
	height   int
	req      request.Model
	tabs     tabs.Model
	base     base.Model
	params   params.Model
	headers  headers.Model
	body     body.Model
	viewport viewport.Model
	mode     helpers.Mode
	err      error
	quitting bool
}

func New() *MainModel {
	return &MainModel{
		req:      *request.New(),
		tabs:     tabs.New(),
		base:     base.New(),
		params:   params.New(),
		headers:  headers.New(),
		body:     body.New(),
		viewport: viewport.New(),
	}
}

func (m *MainModel) makeRequest() {
	m.req.SetURL(m.base.GetURL())
	m.req.SetMethod(m.base.GetMethod())
	m.req.SetParams(m.params.GetParams())
	m.req.SetHeaders(m.headers.GetHeaders())
	body, cType := m.body.GetBody()
	m.req.SetContentType(cType)
	curl := exec.Command("curl", "-X", m.req.GetMethod(), m.req.GetURL())
	for k, v := range m.req.GetHeaders() {
		curl.Args = append(curl.Args, "-H", k+": "+v)
	}
	curl.Args = append(curl.Args, "-d", body)
	output, err := curl.Output()
	if err != nil {
		m.viewport.SetContent("Error running curl check your URL and try again")
		return
	}

	var data interface{}
	json.Unmarshal(output, &data)
	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		m.viewport.SetContent(string(output))
		return
	}
	m.viewport.SetContent(string(formattedJSON))
}

func (m MainModel) renderTab() string {
	var tabContent = map[helpers.Tab]string{
		helpers.Params:  m.params.View(),
		helpers.Auth:    "Authorization Tab\n\nComing Soon\n\nUse the headers tab for now",
		helpers.Headers: m.headers.View(),
		helpers.Body:    m.body.View(),
		helpers.Base:    m.base.View(),
	}

	return styles.Spacer.Render(lipgloss.PlaceVertical(m.height/3, lipgloss.Top, tabContent[m.tabs.GetFocused()]))
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
			m.body.SetWidth(m.width)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.mode.IsInsert() {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			m.mode = helpers.Normal
		case "i":
			m.mode = helpers.Insert
		case "v":
			if !m.mode.IsInsert() {
				m.mode = helpers.Visual
			}
		case "R":
			if !m.mode.IsInsert() {
				m.makeRequest()
			}
		}
	}

	var cmd tea.Cmd
	if !m.mode.IsInsert() {
		m.tabs, cmd = m.tabs.Update(msg)
	}
	if m.mode.IsVisual() {
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	switch m.tabs.GetFocused() {
	case helpers.Params:
		m.params, cmd = m.params.Update(msg)
	case helpers.Body:
		m.body, cmd = m.body.Update(msg)
	case helpers.Headers:
		m.headers, cmd = m.headers.Update(msg)
	default:
		m.base, cmd = m.base.Update(msg)
	}
	return m, cmd
}

func (m MainModel) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		if m.width < 54 {
			return lipgloss.NewStyle().Padding(3, 3).Render("Window Width too small to render app\n\nPlease resize")
		}

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
