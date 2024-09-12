package python

import (
	"go/token"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/event"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonscanner"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythontype"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

func incrementFromDiff(diff *event.Diff) *pythonscanner.Increment {
	var update *pythonscanner.Increment
	switch diff.GetType() {
	case event.DiffType_INSERT:
		update = &pythonscanner.Increment{
			Begin:       token.Pos(diff.GetOffset()),
			End:         token.Pos(diff.GetOffset()),
			Replacement: []byte(diff.GetText()),
		}
	case event.DiffType_DELETE:
		update = &pythonscanner.Increment{
			Begin: token.Pos(diff.GetOffset()),
			End:   token.Pos(diff.GetOffset() + int32(len(diff.GetText()))),
		}
	}
	return update
}

// GetExternalSymbol for the provided value
func GetExternalSymbol(ctx fastnodectx.Context, rm pythonresource.Manager, val pythontype.Value) (pythonresource.Symbol, error) {
	return pythontype.ChooseExternal(ctx, rm, val)
}

// GetExternalSymbols for the provided value
func GetExternalSymbols(ctx fastnodectx.Context, rm pythonresource.Manager, val pythontype.Value) []pythonresource.Symbol {
	return pythontype.AllExternals(ctx, rm, val)
}
