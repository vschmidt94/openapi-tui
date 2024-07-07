package models

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/pb33f/libopenapi/renderer"
	"github.com/vschmidt94/openapi-tui/lib/config"
)

type view int

const (
	sitesList view = iota
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
	siteList := NewSiteList(cfg)
	m.submodels = append(m.submodels, siteList)
	endpointsModel := NewEndpointsModel()
	m.submodels = append(m.submodels, endpointsModel)
	m.activeView = sitesList

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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return nil, tea.Quit
		case "esc":
			if m.activeView == schemaView {
				sl, _ := m.submodels[sitesList].(siteListModel)
				sl.selected = false
				m.submodels[sitesList] = sl
				return m, tea.ClearScreen
			}
		}
	}

	model, cmd = m.submodels[m.activeView].(tea.Model).Update(msg)
	m.submodels[m.activeView] = model
	cmds = append(cmds, cmd)

	if m.submodels[sitesList].(siteListModel).selected {
		nextView = schemaView
	} else {
		nextView = sitesList
	}

	if m.activeView != nextView {
		m.activeView = nextView
		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.activeView {
	case sitesList:
		return m.submodels[sitesList].(siteListModel).View()
	case schemaView:
		return m.submodels[schemaView].(Endpoint).View()
	}
	return "Unknown View"
}
