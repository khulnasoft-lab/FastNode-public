package editorapi

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// Endpoint describes the endpoints that a language must satisfy
// in order to serve the editor API.
// All endpoints return the specified response, a response code, and an error.
type Endpoint interface {
	Language() lang.Language
	ValueReport(ctx fastnodectx.Context, id string) (*ReportResponse, int, error)
	ValueMembersExt(ctx fastnodectx.Context, id string, offset, limit int) (*MembersExtResponse, int, error)
	ValueMembers(ctx fastnodectx.Context, id string, offset, limit int) (*MembersResponse, int, error)
	SymbolReport(ctx fastnodectx.Context, id string) (*ReportResponse, int, error)
	Search(ctx fastnodectx.Context, query string, offset, limit int) (*SearchResults, int, error)
}
