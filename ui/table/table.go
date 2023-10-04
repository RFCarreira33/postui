package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcarreira33/postui/styles"
)

type Model struct {
	table table.Model
}

func New(columns map[string]int, height int) *Model {
	formatedColumns := []table.Column{}
	for k, v := range columns {
		formatedColumns = append(formatedColumns, table.Column{Title: k, Width: v})
	}

	t := table.New(
		table.WithColumns(formatedColumns),
		table.WithHeight(5),
	)

	t.SetStyles(styles.TableStyle)

	return &Model{table: t}
}

// Deletes the selected row from the table and returns the updated rows
func (m *Model) Delete() {
	selectedRow := m.table.SelectedRow()
	rows := m.table.Rows()
	var updatedRows []table.Row
	for i, row := range rows {
		if row[0] == selectedRow[0] && row[1] == selectedRow[1] {
			continue
		}
		updatedRows = append(updatedRows, rows[i])
	}
	m.table.SetRows(updatedRows)
}

func (m *Model) Append(key, value string) {
	m.table.SetRows(append(m.table.Rows(), table.Row{key, value}))
}

// returns the content of the table as a map
func (m Model) Get() map[string]string {
	rows := m.table.Rows()
	params := make(map[string]string)
	for _, row := range rows {
		params[row[0]] = row[1]
	}
	return params
}

func (m *Model) Focus() {
	m.table.Focus()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}
