package pythonanalyzer

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonast"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonstatic"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythontype"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// Models to resolve an ast against.
type Models struct {
	Importer pythonstatic.Importer
	// Shadow module from offline analysis, may be nil.
	Shadow *pythontype.SourceModule
}

// Resolve expressions in the specified module.
func Resolve(ctx fastnodectx.Context, m Models, ast *pythonast.Module, opts Options) (*ResolvedAST, error) {
	ctx.CheckAbort()

	res := NewResolverUsingImporter(m.Importer, opts)

	return res.ResolveContext(ctx, ast, false)
}
