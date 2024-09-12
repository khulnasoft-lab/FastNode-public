package driver

import (
	"net/http"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/core"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// State contains the file driver and buffer handler for a specific buffer state
// associated with a filename and editor.
type State struct {
	Filename      string
	Editor        string
	State         string
	FileDriver    core.FileDriver
	BufferHandler http.Handler
}

// Provider is an interface for querying for driver state
type Provider interface {
	Driver(ctx fastnodectx.Context, filename, editor, state string) (*State, bool)
	DriverFromContent(ctx fastnodectx.Context, filename, editor, content string, cursor int) *State
	LatestDriver(ctx fastnodectx.Context, unixFilepath string) *python.UnifiedDriver
}
