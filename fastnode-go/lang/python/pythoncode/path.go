package pythoncode

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonimports"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythontype"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// ValuePath returns a pythonimports.DottedPath for an external or builtin value. If the path is not found, an empty
// DottedPath is returned.
func ValuePath(ctx fastnodectx.Context, val pythontype.Value, rm pythonresource.Manager) pythonimports.DottedPath {
	if val == nil {
		return pythonimports.DottedPath{}
	}

	val = pythontype.Translate(ctx, pythontype.WidenConstants(val), rm)

	if val == nil {
		return pythonimports.DottedPath{}
	}

	if global, ok := val.(pythontype.GlobalValue); ok {
		switch global := global.(type) {
		case pythontype.ExternalInstance:
			return global.TypeExternal.Symbol().Path()
		case pythontype.External:
			return global.Symbol().Path()
		default:
			return pythonimports.DottedPath{}
		}
	}

	return pythonimports.DottedPath{}
}
