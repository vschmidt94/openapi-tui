package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/vschmidt94/openapi-tui/lib/config"
	"github.com/vschmidt94/openapi-tui/types"
	"io"
	"strings"
)

//type form interface{}

type listState int

const (
	Normal listState = iota
	Selected
	UpdateRequested
)

const NEW_SITE_IDX = -1

type SiteListModel struct {
	Sites        list.Model
	state        listState
	updateForm   *updateForm
	err          error
	windowWidth  int
	windowHeight int
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

func NewSiteList(cfg config.Config) SiteListModel {
	m := SiteListModel{}
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
	m.updateForm = nil

	return m
}

func (m SiteListModel) Init() tea.Cmd {
	m.state = Normal
	return nil
}

func (m SiteListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		s   types.Site
		uf  updateForm
		idx int
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setWindowSize(msg)
	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case "enter":
			fmt.Println("Selected site:", m.Sites.SelectedItem())
			m.state = Selected
		case "u", "n":
			if k == "u" {
				idx = m.Sites.Cursor()
				s = m.SelectedSite()
			} else {
				idx = NEW_SITE_IDX
				s = types.Site{
					Name: "New Site",
					Uri:  "https://example.com",
				}
			}
			uf = NewUpdateForm(&s, idx)
			m.state = UpdateRequested
			m.updateForm = &uf
		}
	}

	m.Sites, _ = m.Sites.Update(msg)
	return m, nil
}

func (m SiteListModel) View() string {
	return m.Sites.View()
}

func (m *SiteListModel) setWindowSize(msg tea.WindowSizeMsg) {
	m.windowHeight = msg.Height
	m.windowWidth = msg.Width
	m.Sites.SetWidth(msg.Width / 3)
	m.Sites.SetHeight(msg.Height - 2)
}

func (m *SiteListModel) SelectedSite() types.Site {
	return m.Sites.SelectedItem().(types.Site)
}
