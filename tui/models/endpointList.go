package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pb33f/libopenapi"
	"github.com/vschmidt94/openapi-tui/types"
	"io"
	"net/http"
)

type Endpoint struct {
	Doc    libopenapi.Document
	Loaded bool
}

func NewEndpointsModel() Endpoint {
	return Endpoint{}
}

func (m Endpoint) Init() tea.Cmd {
	return nil
}

func (m Endpoint) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgOpenApiDocRequest:
		cmd := getOpenApiSchema(msg.Site)
		return m, cmd
	case MsgOpenApiDocResponse:
		m.Doc = msg.doc
		m.Loaded = true
	}

	return m, nil
}

func (m Endpoint) View() string {
	if !m.Loaded {
		return "Loading..."
	}

	return fmt.Sprintf("%v", m.Doc)
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
			return fmt.Errorf("failed to parse into OpenAPI Doc: %v", err)
		}

		return MsgOpenApiDocResponse{
			doc:  doc,
			site: site,
		}
	}
}
