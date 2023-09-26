package body

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/app/helpers"
	"github.com/rfcarreira33/postui/styles"
	"github.com/rfcarreira33/postui/ui/selector"
)

type Model struct {
	focusedInput int
	inputs       []textinput.Model
	textarea     textarea.Model
	mode         helpers.Mode
	paramsTable  table.Model
	selector     selector.Model
	err          error
}

func New() Model {
	// initialize textarea
	ta := textarea.New()
	ta.Placeholder = "Body Data"

	// initialize text inputs and the rest of the model
	m := Model{
		inputs:   make([]textinput.Model, 2),
		textarea: ta,
		selector: *selector.New([]string{"JSON", "Javascript", "HTML", "XML", "Text", "FormData"}),
	}

	var input textinput.Model
	for i := range m.inputs {
		input = textinput.New()

		switch i {
		case 0:
			input.Placeholder = "Enter Key"
		case 1:
			input.Placeholder = "Enter Value"
		}
		m.inputs[i] = input
	}

	columns := []table.Column{
		{Title: "Key", Width: 20},
		{Title: "Value", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(5),
	)

	t.SetStyles(styles.TableStyle)

	m.paramsTable = t

	return m
}

func (m *Model) SetWidth(width int) {
	m.textarea.SetWidth(width)
}

func (m *Model) cycleFocus() {
	m.inputs[m.focusedInput].Blur()
	if m.focusedInput == len(m.inputs)-1 {
		m.focusedInput = 0
	} else {
		m.focusedInput++
	}
	m.inputs[m.focusedInput].Focus()
}

func (m *Model) appendData() {
	key := m.inputs[0].Value()
	val := m.inputs[1].Value()
	if key == "" || val == "" {
		return
	}
	m.inputs[0].SetValue("")
	m.inputs[1].SetValue("")
	m.inputs[m.focusedInput].Blur()
	m.focusedInput = 0
	m.inputs[m.focusedInput].Focus()
	m.paramsTable.SetRows(append(m.paramsTable.Rows(), table.Row{key, val}))
}

func (m Model) GetFormData() map[string]string {
	rows := m.paramsTable.Rows()
	params := make(map[string]string)
	for _, row := range rows {
		params[row[0]] = row[1]
	}
	return params
}

// return body, content-type
func (m Model) GetBody() (string, string) {
	var body string
	cType := m.selector.GetSelectedItem()
	if cType == "FormData" {
		for k, v := range m.GetFormData() {
			body += k + "=" + v + "&"
		}
		return body, cType
	}
	return m.textarea.Value(), cType
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	formData := m.selector.GetSelectedItem() == "FormData"

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = helpers.Normal
			m.inputs[m.focusedInput].Blur()
			m.textarea.Blur()
		case "i":
			if m.mode.IsNormal() {
				m.mode = helpers.Insert
				if formData {
					m.inputs[m.focusedInput].Focus()
				} else {
					m.textarea.Focus()
				}
				return m, nil
			}
		case "tab":
			if m.mode.IsInsert() && formData {
				m.cycleFocus()
			}
		case "enter":
			if formData {
				m.appendData()
			}
		case "d":
			if !m.mode.IsInsert() {
				selectedRow := m.paramsTable.SelectedRow()
				rows := m.paramsTable.Rows()
				var updatedRows []table.Row
				for i, row := range rows {
					if row[0] == selectedRow[0] && row[1] == selectedRow[1] {
						continue
					}
					updatedRows = append(updatedRows, rows[i])
				}
				m.paramsTable.SetRows(updatedRows)
			}
		case "P":
			if !m.mode.IsInsert() && formData {
				m.mode = helpers.Visual
				m.paramsTable.Focus()
			}
		}
	}

	if m.mode.IsInsert() {
		if formData {
			m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
		} else {
			m.textarea, cmd = m.textarea.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	if m.mode.IsNormal() {
		m.selector, cmd = m.selector.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.paramsTable, cmd = m.paramsTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	inputs := lipgloss.JoinVertical(lipgloss.Top, "Key\t"+m.inputs[0].View(), "\nValue \t"+m.inputs[1].View(), "\n\n Check and or remove added params with 'P'")
	table := m.paramsTable.View() + "\n Press 'd' to delete selected row"

	body := m.textarea.View()
	if m.selector.GetSelectedItem() == "FormData" {
		if m.mode.IsVisual() {
			return table
		}
		body = inputs
	}

	return lipgloss.JoinVertical(lipgloss.Top, m.selector.View(), "\n"+body)
}
