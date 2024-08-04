package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/vschmidt94/openapi-tui/types"

	// "github.com/pb33f/libopenapi/renderer"
	"github.com/vschmidt94/openapi-tui/lib/config"
)

type view int

const (
	sitesListView view = iota
	updateSiteView
	schemaView
)

type submodel interface{}

type mainModel struct {
	cfg        config.Config
	submodels  []submodel
	activeView view
}

func New(cfg config.Config) *mainModel {
	m := &mainModel{}
	m.cfg = cfg
	// order matters to align with view enum
	siteList := NewSiteList(cfg)
	m.submodels = append(m.submodels, siteList)
	m.submodels = append(m.submodels, nil)
	endpointsModel := NewEndpointsModel()
	m.submodels = append(m.submodels, endpointsModel)

	m.activeView = sitesListView

	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd      tea.Cmd
		cmds     []tea.Cmd
		model    tea.Model
		nextView view
	)

	// handle global key events
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return nil, tea.Quit
		case "esc":
			if m.activeView == schemaView || m.activeView == updateSiteView {
				m.activeView = sitesListView
				sl, _ := m.submodels[sitesListView].(SiteListModel)
				sl.state = Normal
				m.submodels[sitesListView] = sl
				return m, tea.ClearScreen
			}
		}
	}

	// send the message to the active view
	model, cmd = m.submodels[m.activeView].(tea.Model).Update(msg)
	m.submodels[m.activeView] = model
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// switch views if needed
	switch m.activeView {
	case sitesListView:
		sl, _ := m.submodels[sitesListView].(SiteListModel)
		switch sl.state {
		case UpdateRequested:
			if sl.updateForm == nil {
				// expected to have a form to use
				sl.state = Normal
				nextView = sitesListView
			} else {
				m.submodels[updateSiteView] = *sl.updateForm
				nextView = updateSiteView
			}
			// reset site list state so we can return to it eventually
			sl.state = Normal
			m.submodels[sitesListView] = sl
		case Selected:
			ep := NewEndpointsModel()
			selectedSite := sl.Sites.SelectedItem().(types.Site)
			ep.site = selectedSite
			m.submodels[schemaView] = ep
			cmd = m.submodels[schemaView].(Endpoint).Init()
			cmds = append(cmds, cmd)
			nextView = schemaView
		}
	case updateSiteView:
		uf, _ := m.submodels[updateSiteView].(updateForm)
		if uf.Form.State == huh.StateCompleted || uf.Form.State == huh.StateAborted {
			// form is done, update the site and return to list
			sl, _ := m.submodels[sitesListView].(SiteListModel)
			if uf.idx == NEW_SITE_IDX {
				sl.Sites.InsertItem(0, *uf.site)
			} else {
				sl.Sites.SetItem(uf.idx, *uf.site)
			}
			m.submodels[sitesListView] = sl
			m.submodels[updateSiteView] = nil
			nextView = sitesListView
		} else {
			// normal state
			nextView = updateSiteView
		}
	case schemaView:
		nextView = schemaView
	}

	if m.activeView != nextView {
		m.activeView = nextView
		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.activeView {
	case sitesListView:
		return m.submodels[sitesListView].(SiteListModel).View()
	case updateSiteView:
		ufp := m.submodels[updateSiteView].(updateForm)
		return ufp.View()
	case schemaView:
		return m.submodels[schemaView].(Endpoint).View()
	}
	return "Unknown View"
}
