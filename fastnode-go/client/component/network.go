package component

import (
	"context"
)

//NetworkManager defines the functions to query whether or not network connectivity exists
type NetworkManager interface {
	Core
	// Online returns whether or not there is network connectivity
	Online() bool

	// CheckOnline checks and returns whether or not there is network connectivity
	CheckOnline(ctx context.Context) bool

	// FastnodedOnline checks and returns whether or not fastnoded has been initialized to
	// an extent where it can reliably report state to requests from clients
	FastnodedOnline() bool

	SetOffline(bool)
}
