package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/styles"
	"github.com/rfcarreira33/postui/ui"
)

type status int

const divisor = 5

const (
	method status = iota
	url
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
	if m.focused == url {
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
		m.focused = url
	} else {
		m.focused--
	}
}

func (m *MainModel) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, ui.ItemDelegate{}, width/divisor*3, height/divisor*2)
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
	m.urlInput.Width = width / divisor * 3
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.initTextInput(msg.Width)
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
			curl := exec.Command("curl", "-X", m.req.method, m.req.url)
			curl.Stdout = os.Stdout
			curl.Stderr = os.Stderr
			curl.Run()
		case "enter":
			if m.focused == url {
				if !m.editing {
					m.urlInput.Focus()
					m.editing = true
				} else {
					m.req.url = m.urlInput.Value()
					m.editing = false
					m.urlInput.Blur()
				}
			}
		}
	}
	var cmd tea.Cmd
	m.urlInput, cmd = m.urlInput.Update(msg)
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		methodView := m.lists[method].View()
		urlView := m.urlInput.View()

		switch m.focused {
		case url:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.ColumnStyle.Render(methodView),
				styles.FocusedStyle.Render(urlView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.FocusedStyle.Render(methodView),
				styles.ColumnStyle.Render(urlView),
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
