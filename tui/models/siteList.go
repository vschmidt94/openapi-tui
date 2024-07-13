package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	// "github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/vschmidt94/openapi-tui/lib/config"
	"github.com/vschmidt94/openapi-tui/types"
	"io"
	"strings"
)

type form interface{}

type siteListModel struct {
	Sites          list.Model
	selected       bool
	err            error
	showUpdateForm bool
	updateForm     form
	windowWidth    int
	windowHeight   int
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
	// var cmd  tea.Cmd
	// var cmds []tea.Cmd

	if m.showUpdateForm {
		uf := m.updateForm.(updateForm)
		switch formState := uf.Form.State; formState {
		case huh.StateNormal:
			m.updateForm, _ = m.updateForm.(tea.Model).Update(msg)
		case huh.StateCompleted:
			m.showUpdateForm = false
			m.updateForm = nil
			if uf.isNew {
				newSite := types.Site{
					Name:         uf.site.Name,
					Uri:          uf.site.Uri,
					User:         uf.site.User,
					RequiresAuth: uf.site.RequiresAuth,
				}
				m.Sites.InsertItem(0, newSite)
			} else {
				idx := m.Sites.Cursor()
				m.Sites.SetItem(idx, *uf.site)
			}
		case huh.StateAborted:
			m.showUpdateForm = false
			m.updateForm = nil
			// TODO: handle abort
		}

		return m, nil
	}

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
		case "u":
			s := m.SelectedSite()
			m.updateForm = NewUpdateForm(&s, false)
			m.showUpdateForm = true
			if uf, ok := m.updateForm.(updateForm); ok {
				uf.Form.Run()
			}
		case "n":
			newSite := types.Site{
				Name: "New Site",
				Uri:  "http://example.com",
			}
			m.updateForm = NewUpdateForm(&newSite, true)
			m.showUpdateForm = true
			if uf, ok := m.updateForm.(updateForm); ok {
				uf.Form.Run()
			}
		}
	}

	m.Sites, _ = m.Sites.Update(msg)
	return m, nil
}

func (m siteListModel) View() string {
	if m.showUpdateForm {
		rendered := lipgloss.JoinHorizontal(lipgloss.Top, m.Sites.View(), m.updateForm.(tea.Model).View())
		return rendered
	}
	return m.Sites.View()
}

func (m *siteListModel) setWindowSize(msg tea.WindowSizeMsg) {
	m.windowHeight = msg.Height
	m.windowWidth = msg.Width
	m.Sites.SetWidth(msg.Width / 3)
	m.Sites.SetHeight(msg.Height - 2)
}

func (m *siteListModel) SelectedSite() types.Site {
	return m.Sites.SelectedItem().(types.Site)
}
