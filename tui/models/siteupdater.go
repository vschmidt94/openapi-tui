package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/vschmidt94/openapi-tui/types"
)

type updateForm struct {
	Form      *huh.Form
	site      *types.Site
	isNew     bool
	isRunning bool
}

func NewUpdateForm(site *types.Site, isNew bool) updateForm {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&site.Name).Key(site.Name),
			huh.NewInput().Title("Uri").Value(&site.Uri).Key(site.Uri),
			huh.NewInput().Title("User").Value(&site.User).Key(site.User),
			huh.NewConfirm().Title("Requires Auth").Value(&site.RequiresAuth).Key("requiresAuth"),
		),
	).
		WithWidth(80).
		WithHeight(25)

	return updateForm{
		isNew: isNew,
		site:  site,
		Form:  form,
	}
}

func (m updateForm) Init() tea.Cmd {
	cmd := m.Form.Init()
	return cmd
}

func (m updateForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := m.Form.Update(msg)
	cmds = append(cmds, cmd)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m updateForm) Run() {
	m.Form.Run()
}

func (m updateForm) View() string {
	//var title string
	//if m.isNew {
	//	title = "Create New API Site"
	//} else {
	//	title = "Update API Site"
	//}
	// titleStyle := lipgloss.NewStyle().Bold(true).Padding(1).Foreground(lipgloss.Color("170"))
	// renderedTitle := titleStyle.Render(title)
	// return lipgloss.JoinVertical(lipgloss.Left, renderedTitle, m.Form.View())
	return m.Form.View()
}
