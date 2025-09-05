// Package list provides a table model for displaying a list of items with title and date columns.
package list

import (
	"fmt"
	"strings"

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

var (
	helpKeyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e738d"))
	helpDescStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#5b6078"))
	helpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#494d64"))
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

	hm := help.New()
	hm.Styles = help.Styles{
		ShortKey:       helpKeyStyle,
		ShortDesc:      helpDescStyle,
		ShortSeparator: helpSeparatorStyle,
		FullKey:        helpKeyStyle,
		FullDesc:       helpDescStyle,
		FullSeparator:  helpSeparatorStyle,
		Ellipsis:       helpSeparatorStyle,
	}

	return Model{
		table: table.
			New(columns).
			Filtered(true).
			Focused(true).
			SelectableRows(true).
			WithKeyMap(km.KeyMap).
			BorderRounded().
			WithBaseStyle(baseStyle).
			HeaderStyle(headerStyle).
			WithRows(rows),
		help:  hm,
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
		row := table.NewRow(table.RowData{
			columnKeyTitle: n.Title,
			columnKeyDate:  n.Date.Format("02 Jan 2006"),
			columnKeyLink:  n.Link,
		})

		for word, style := range keywordStyles {
			if strings.Contains(strings.ToLower(n.Title), word) {
				row = row.WithStyle(style)
				break
			}
		}

		rows = append(rows, row)
	}

	return rows
}

var keywordStyles = map[string]lipgloss.Style{
	"exam":         lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true),
	"registration": lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	"payment":      lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
	"make up":      lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	"holiday":      lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width := min(msg.Width, m.width)
		height := msg.Height
		m.help.Width = width
		m.table = m.table.
			WithTargetWidth(width).
			WithPageSize(height - 5)
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
			// Reset filter and selection after opening
			m.table = m.table.
				WithFilterInputValue("").
				WithAllRowsDeselected()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			if !m.table.GetIsFilterInputFocused() {
				cmds = append(cmds, tea.Quit)
				fmt.Println()
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	tableView := m.table.View()
	helpView := m.help.View(m.keys)
	return lipgloss.JoinVertical(lipgloss.Left, tableView, helpView)
}
