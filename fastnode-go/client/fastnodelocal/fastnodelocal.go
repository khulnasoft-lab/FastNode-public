// Package fastnodelocal re-exports useful pieces of the internal package fastnodelocal for testing, etc
package fastnodelocal

import (
	"context"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
)

// LoadOptions specifies various loading options for python services
type LoadOptions = fastnodelocal.LoadOptions

// LoadResourceManager loads the resource manager for Fastnode Local
func LoadResourceManager(ctx context.Context, opts LoadOptions) (pythonresource.Manager, error) {
	return fastnodelocal.LoadResourceManager(ctx, opts)
}
