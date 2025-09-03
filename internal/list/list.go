// Package list provides a table model for displaying a list of items with title and date columns.
package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/AtifChy/aiub-notice/internal/notice"
)

const (
	columnKeyTitle = "title"
	columnKeyDate  = "date"
	columnKeyLink  = "link"
)

var (
	baseStyle   = lipgloss.NewStyle().BorderForeground(lipgloss.Color("0"))
	headerStyle = lipgloss.NewStyle().Align(lipgloss.Center).Bold(true).Foreground(lipgloss.Color("6"))
)

type Model struct {
	table table.Model
	keys  KeyMap
	help  help.Model
	width int
}

func NewModel() Model {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyTitle, "Title", 6).
			WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)).
			WithFiltered(true),
		table.NewFlexColumn(columnKeyDate, "Date", 2).
			WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
	}
	rows := getRows()

	km := DefaultKeyMap()

	return Model{
		table: table.
			New(columns).
			Filtered(true).
			WithFuzzyFilter().
			Focused(true).
			WithPageSize(10).
			SelectableRows(true).
			WithKeyMap(km.KeyMap).
			BorderRounded().
			WithBaseStyle(baseStyle).
			HeaderStyle(headerStyle).
			WithRows(rows),
		help:  help.New(),
		keys:  km,
		width: 90,
	}
}

func getRows() []table.Row {
	notices, err := notice.LoadNoticesCache()
	if err != nil {
		fmt.Println("Error loading notices from cache:", err)
		return []table.Row{}
	}

	var rows []table.Row

	for _, n := range notices {
		tableRow := table.NewRow(table.RowData{
			columnKeyTitle: n.Title,
			columnKeyDate:  n.Date.Format("02 Jan 2006"),
			columnKeyLink:  n.Link,
		})
		rows = append(rows, tableRow)
	}

	return rows
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width := msg.Width
		height := msg.Height

		if width > m.width {
			width = m.width
		}

		m.help.Width = m.width

		m.table = m.table.WithTargetWidth(width).WithPageSize(height - 5) // adjust rows to fit height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.RowOpen):
			if m.table.GetIsFilterInputFocused() {
				break
			}
			rows := m.table.SelectedRows()
			if len(rows) == 0 {
				rows = []table.Row{m.table.HighlightedRow()}
			}
			for _, row := range rows {
				if val, ok := row.Data[columnKeyLink]; ok && val != nil {
					link := val.(string)
					if err := openURL(link); err != nil {
						fmt.Println("Error opening URL:", err)
					}
				}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			if !m.table.GetIsFilterInputFocused() {
				cmds = append(cmds, tea.Quit)
				fmt.Println()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	tableView := m.table.View()
	helpView := m.help.View(m.keys)
	return lipgloss.JoinVertical(lipgloss.Left, tableView, helpView)
}
