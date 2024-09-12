package pythonresource

import (
	"io"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/editorapi"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythoncode/symbolcounts"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonimports"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource/internal/resources/docs"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource/keytypes"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// Manager is the Python resource manager
type Manager interface {
	io.Closer
	Reset()

	Distributions() []keytypes.Distribution
	DistLoaded(dist keytypes.Distribution) bool

	ArgSpec(sym Symbol) *pythonimports.ArgSpec
	PopularSignatures(sym Symbol) []*editorapi.Signature
	CumulativeNumArgsFrequency(sym Symbol, numArgs int) (float64, bool)
	KeywordArgFrequency(sym Symbol, arg string) (int, bool)
	NumArgsFrequency(sym Symbol, numArgs int) (float64, bool)
	Documentation(sym Symbol) *docs.Entity
	SymbolCounts(sym Symbol) *symbolcounts.Counts
	Kwargs(sym Symbol) *KeywordArgs
	TruthyReturnTypes(sym Symbol) []TruthySymbol
	ReturnTypes(sym Symbol) []Symbol

	PathSymbol(path pythonimports.DottedPath) (Symbol, error)
	PathSymbols(ctx fastnodectx.Context, path pythonimports.DottedPath) ([]Symbol, error)
	NewSymbol(dist keytypes.Distribution, path pythonimports.DottedPath) (Symbol, error)
	Kind(s Symbol) keytypes.Kind
	Type(s Symbol) (Symbol, error)
	Bases(s Symbol) []Symbol
	Children(s Symbol) ([]string, error)
	ChildSymbol(s Symbol, c string) (Symbol, error)
	CanonicalSymbols(dist keytypes.Distribution) ([]Symbol, error)
	TopLevels(dist keytypes.Distribution) ([]string, error)

	Pkgs() []string
	DistsForPkg(pkg string) []keytypes.Distribution

	SigStats(sym Symbol) *SigStats
}
