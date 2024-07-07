package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vschmidt94/openapi-tui/lib/config"
	"github.com/vschmidt94/openapi-tui/types"
	"io"
	"strings"
)

type siteListModel struct {
	Sites    list.Model
	selected bool
	err      error
}

/* List Item Delegate for styling */
var (
	listTitleStyle    = lipgloss.NewStyle().MarginLeft(2).Underline(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, idx int, listItem list.Item) {
	i, ok := listItem.(types.Site)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Name)

	fn := itemStyle.Render
	if idx == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
func NewSiteList(cfg config.Config) siteListModel {
	m := siteListModel{}
	m.Sites = list.New([]list.Item{}, itemDelegate{}, 0, 0)
	m.Sites.Title = "API Sites"
	m.Sites.SetShowStatusBar(false)
	m.Sites.Styles.Title = listTitleStyle
	m.Sites.Styles.HelpStyle = helpStyle
	m.Sites.Styles.PaginationStyle = paginationStyle
	var sites []list.Item
	for _, site := range cfg.Sites {
		sites = append(sites, site)
	}
	m.Sites.SetItems(sites)

	return m
}

func (m siteListModel) Init() tea.Cmd {
	return nil
}

func (m siteListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setWindowSize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			fmt.Println("Selected site:", m.Sites.SelectedItem())
			m.selected = true
		}
	}

	m.Sites, _ = m.Sites.Update(msg)
	return m, nil
}

func (m siteListModel) View() string {
	return m.Sites.View()
}

func (m *siteListModel) setWindowSize(msg tea.WindowSizeMsg) {
	m.Sites.SetWidth(msg.Width - 2)
	m.Sites.SetHeight(msg.Height - 2)
}

func (m *siteListModel) SelectedSite() types.Site {
	return m.Sites.SelectedItem().(types.Site)
}
