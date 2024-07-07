package models

import (
	"github.com/pb33f/libopenapi"
	"github.com/vschmidt94/openapi-tui/types"
)

type MsgOpenApiDocRequest struct {
	Site types.Site
}

type MsgOpenApiDocResponse struct {
	doc  libopenapi.Document
	site types.Site
}
