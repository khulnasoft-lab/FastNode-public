package pythonstatic

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonast"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// AnalyzeGlobal analyzes source using only the global graph, in particular this means that
// the address assigned to symbols defined in the module is ARBITRARY.
func AnalyzeGlobal(ctx fastnodectx.Context, ai AssemblerInputs, ast *pythonast.Module) (*Assembly, error) {
	opts := DefaultOptions
	opts.AllowValueMutation = true
	assembler := NewAssembler(ctx, ai, opts)
	assembler.AddSource(ASTBundle{
		AST:     ast,
		Path:    "/src.py",
		Imports: FindImports(ctx, "/src.py", ast),
	})
	return assembler.Build(ctx)
}
