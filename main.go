package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/styles"
	"github.com/rfcarreira33/postui/ui"
)

type status int

const Xdivisor = 4
const Ydivisor = 3

const (
	method status = iota
	url
	viewer
)

type Request struct {
	url    string
	method string
}

type MainModel struct {
	loaded   bool
	focused  status
	urlInput textinput.Model
	lists    []list.Model
	viewport viewport.Model
	help     help.Model
	err      error
	quitting bool
	editing  bool
	req      Request
}

func New() *MainModel {
	return &MainModel{
		req: Request{
			url:    "",
			method: "GET",
		},
	}
}

// Move focus
func (m *MainModel) Next() {
	if m.editing {
		return
	}
	if m.focused == viewer {
		m.focused = method
	} else {
		m.focused++
	}
}

func (m *MainModel) Prev() {
	if m.editing {
		return
	}
	if m.focused == method {
		m.focused = viewer
	} else {
		m.focused--
	}
}

// Init Functions
func (m *MainModel) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, ui.ItemDelegate{}, width/Xdivisor, height/Ydivisor)
	defaultList.SetShowHelp(false)
	defaultList.SetShowStatusBar(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[method].Title = "HTTP Method"
	m.lists[method].SetItems([]list.Item{
		ui.Item("GET"),
		ui.Item("POST"),
		ui.Item("PUT"),
		ui.Item("PATCH"),
		ui.Item("DELETE"),
	})
	m.lists[method].SetShowTitle(false)
	m.lists[method].SetFilteringEnabled(false)
}

func (m *MainModel) initTextInput(width int) {
	m.urlInput = textinput.New()
	m.urlInput.Placeholder = "URL"
	m.urlInput.CharLimit = 100
	m.urlInput.Width = width / Xdivisor * 2
}

func (m *MainModel) initViewPort(width, height int) {
	m.viewport = viewport.New(width, height/Ydivisor-Ydivisor)
}

func (m *MainModel) initHelp(width int) {
	m.help = help.New()
	m.help.Width = width / Xdivisor
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle key presses
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.initTextInput(msg.Width)
			m.initViewPort(msg.Width, msg.Height)
			m.initHelp(msg.Width)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "h", "left":
			m.Prev()
		case "l", "right":
			m.Next()
		case "R":
			m.req.method = string(m.lists[method].SelectedItem().(ui.Item))
			curl := exec.Command("curl", "-X", m.req.method, "-H", "Content-Type: application/json", m.req.url)
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
		case "?":
			if !m.editing {
				m.help.ShowAll = !m.help.ShowAll
			}
		case "enter":
			if m.focused == url {
				if !m.editing {
					m.urlInput.Focus()
					m.editing = !m.editing
				} else {
					m.req.url = m.urlInput.Value()
					m.editing = !m.editing
					m.urlInput.Blur()
				}
			}
		}
	}

	// Update the focused component
	var cmd tea.Cmd
	switch m.focused {
	case method:
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	case url:
		m.urlInput, cmd = m.urlInput.Update(msg)
	case viewer:
		m.viewport, cmd = m.viewport.Update(msg)
	}
	m.help, _ = m.help.Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		methodView := m.lists[method].View()
		urlView := m.urlInput.View()
		viewerView := m.viewport.View()
		helpView := m.help.View(ui.Keys)

		switch m.focused {
		case url:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.JoinVertical(
					lipgloss.Top,
					styles.ColumnStyle.Render(methodView),
					styles.ColumnStyle.Render(helpView),
				),
				lipgloss.JoinVertical(
					lipgloss.Top,
					styles.FocusedStyle.Render(urlView),
					styles.ColumnStyle.Render(viewerView),
				),
			)
		case viewer:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.JoinVertical(
					lipgloss.Top,
					styles.ColumnStyle.Render(methodView),
					styles.ColumnStyle.Render(helpView),
				),
				lipgloss.JoinVertical(
					lipgloss.Top,
					styles.ColumnStyle.Render(urlView),
					styles.FocusedStyle.Render(viewerView),
				),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.JoinVertical(
					lipgloss.Top,
					styles.FocusedStyle.Render(methodView),
					styles.ColumnStyle.Render(helpView),
				),
				lipgloss.JoinVertical(
					lipgloss.Top,
					styles.ColumnStyle.Render(urlView),
					styles.ColumnStyle.Render(viewerView),
				),
			)
		}
	} else {
		return "loading..."
	}
}

func main() {
	mainModel := New()
	program := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
