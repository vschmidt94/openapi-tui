package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pb33f/libopenapi"
	"github.com/vschmidt94/openapi-tui/types"
	"io"
	"net/http"
	"sort"
	"strings"
)

type state int

const (
	docNew state = iota
	docLoading
	docLoaded
	endpointsReady
)

type epItemDelegate struct{}

func (d epItemDelegate) Height() int                             { return 1 }
func (d epItemDelegate) Spacing() int                            { return 0 }
func (d epItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d epItemDelegate) Render(w io.Writer, m list.Model, idx int, listItem list.Item) {
	i, ok := listItem.(types.Endpoint)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Path)

	fn := itemStyle.Render
	if idx == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type EndpointListModel struct {
	site         types.Site
	Endpoints    list.Model
	currState    state
	doc          libopenapi.Document
	windowHeight int
	windowWidth  int
}

func (m *EndpointListModel) setWindowSize(msg tea.WindowSizeMsg) {
	m.windowHeight = msg.Height
	m.windowWidth = msg.Width
	m.Endpoints.SetWidth(msg.Width / 2)
	m.Endpoints.SetHeight(msg.Height - 2)
}

func NewEndpointsModel() EndpointListModel {
	m := EndpointListModel{}
	m.Endpoints = list.New([]list.Item{}, epItemDelegate{}, 0, 0)
	m.Endpoints.Title = "Endpoints"
	m.Endpoints.SetShowStatusBar(false)
	m.Endpoints.Styles.Title = listTitleStyle
	m.Endpoints.Styles.HelpStyle = helpStyle
	m.Endpoints.Styles.PaginationStyle = paginationStyle
	return EndpointListModel{
		currState: docNew,
	}
}

func (m *EndpointListModel) PopulateEndpoints() {
	if m.doc == nil {
		panic("document is not initialized")
	}
	docModel, errors := m.doc.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("connot create v3 model from document: %d errors reported", len(errors)))
	}
	var endpoints []list.Item
	epl := list.New([]list.Item{}, epItemDelegate{}, 40, 40)
	for pathPair := docModel.Model.Paths.PathItems.First(); pathPair != nil; pathPair = pathPair.Next() {
		pathName := pathPair.Key
		ep := types.Endpoint{}
		ep.Path = pathName()
		ep.Method = "foo"
		ep.Description = "bar"
		endpoints = append(endpoints, ep)
	}
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].(types.Endpoint).Path < endpoints[j].(types.Endpoint).Path
	})
	epl.SetItems(endpoints)
	m.Endpoints = epl
	m.currState = endpointsReady
}

func (m EndpointListModel) Init() tea.Cmd {
	return getOpenApiSchema(m.site)
}

func (m EndpointListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	cmd = nil
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setWindowSize(msg)
	case MsgOpenApiDocRequest:
		cmd = getOpenApiSchema(msg.Site)
		m.currState = docLoading
	case MsgOpenApiDocResponse:
		m.doc = msg.doc
		m.currState = docLoaded
		m.PopulateEndpoints()
	}
	if m.currState == endpointsReady {
		m.Endpoints, cmd = m.Endpoints.Update(msg)
	}

	return m, cmd
}

func (m EndpointListModel) View() string {
	switch m.currState {
	case docNew:
		return "Loading..."
	case docLoading:
		return "Loading..."
	case endpointsReady:
		return m.Endpoints.View()
	}
	return "Unknown state"
}

func getOpenApiSchema(site types.Site) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(site.Uri)
		if err != nil {
			return fmt.Errorf("failed to get schema: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		doc, err := libopenapi.NewDocument(body)
		if err != nil {
			panic(fmt.Sprintf("cannot create new document: %e", err))
		}

		return MsgOpenApiDocResponse{
			doc:  doc,
			site: site,
		}
	}
}
