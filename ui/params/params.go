package params

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcarreira33/postui/app/helpers"
	"github.com/rfcarreira33/postui/styles"
)

type Model struct {
	focusedInput int
	inputs       []textinput.Model
	mode         helpers.Mode
	paramsTable  table.Model
	err          error
}

func New() Model {
	m := Model{
		inputs: make([]textinput.Model, 2),
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

func (m *Model) cycleFocus() {
	m.inputs[m.focusedInput].Blur()
	if m.focusedInput == len(m.inputs)-1 {
		m.focusedInput = 0
	} else {
		m.focusedInput++
	}
	m.inputs[m.focusedInput].Focus()
}

func (m *Model) appendParam() {
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

func (m Model) GetParams() map[string]string {
	rows := m.paramsTable.Rows()
	params := make(map[string]string)
	for _, row := range rows {
		params[row[0]] = row[1]
	}
	return params
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
			m.inputs[m.focusedInput].Blur()
		case "i":
			if m.mode.IsNormal() {
				m.mode = helpers.Insert
				m.inputs[m.focusedInput].Focus()
				return m, nil
			}
		case "tab":
			if m.mode.IsInsert() {
				m.cycleFocus()
			}
		case "enter":
			m.appendParam()
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
			if !m.mode.IsInsert() {
				m.mode = helpers.Visual
				m.paramsTable.Focus()
			}
		}
	}

	if m.mode.IsInsert() {
		m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
		cmds = append(cmds, cmd)
	}
	m.paramsTable, cmd = m.paramsTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	inputs := lipgloss.JoinVertical(lipgloss.Top, "\t Key\t"+m.inputs[0].View(), "\n\t Value \t"+m.inputs[1].View(), "\n\n\t Check and or remove added params with 'P'")
	table := m.paramsTable.View() + "\n\n Press 'd' to delete selected row"

	if m.mode.IsVisual() {
		return table
	}

	return inputs
}