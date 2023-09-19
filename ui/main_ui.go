package ui

import (
	"encoding/json"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	styles "github.com/rfcarreira33/postui/styles"
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
		url:     "https://dummyjson.com/products/1",
		method:  "GET",
		headers: headers,
	}

	return &MainModel{
		req:      req,
		tabs:     tabs.New(),
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
			m.width = msg.Width - 10
			m.height = msg.Height
			m.tabs.SetWidth(msg.Width)
			m.viewport.SetSize(msg.Width, msg.Height)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "v":
			m.view = !m.view
		case "i":
			m.insert = !m.insert
		case "R":
			m.req.method = "GET"
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

	// Update the focused component
	var cmd tea.Cmd
	if m.view {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	m.tabs, cmd = m.tabs.Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		sb := strings.Builder{}
		sb.WriteString(m.tabs.View())
		sb.WriteString("\n")
		sb.WriteString(lg.Place(m.width, m.height/3+4, lg.Left, lg.Top, "Press 'R' to run curl command"))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", m.width+4))
		sb.WriteString("\n")
		sb.WriteString(m.viewport.View())
		return styles.FocusedStyle.Render(sb.String())
	} else {
		return "loading..."
	}
}
