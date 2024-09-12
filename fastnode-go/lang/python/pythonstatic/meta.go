package pythonstatic

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythontype"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

var metaclasses []metaclass

func registerMetaClass(mc metaclass) {
	metaclasses = append(metaclasses, mc)
}

type metaclass interface {
	Construct(ctx fastnodectx.Context, c *pythontype.SourceClass, symbols *pythontype.SymbolTable)
	IsMetaClass(c *pythontype.SourceClass) bool
}

func updateClass(ctx fastnodectx.Context, c *pythontype.SourceClass, symbols *pythontype.SymbolTable) {
	ctx.CheckAbort()

	for _, mc := range metaclasses {
		if mc.IsMetaClass(c) {
			mc.Construct(ctx, c, symbols)
			return
		}
	}

	// default: merge symbol tables
	for name, symbol := range symbols.Table {
		member := c.Members.LocalOrCreate(name)
		member.Value = pythontype.Unite(ctx, member.Value, symbol.Value)
		if !symbol.Name.Equals(member.Name) {
			panic("different address for symbols with the same name")
		}
	}
}

func doParameterHeuristics(ctx fastnodectx.Context, f *pythontype.SourceFunction, prop *propagator) {
	ctx.CheckAbort()

	switch {
	case isDjangoAdmin(prop.Module):
		if p := djangoAdminRequestParam(f); isDjangoRequestParam(p) {
			updateSymbolWithDjangoRequest(ctx, p.Symbol, prop.Importer)
		}
		if p := djangoAdminQuerySetParam(f); isDjangoQuerySetParam(p) {
			updateSymbolWithDjangoQuerySet(ctx, p.Symbol, f)
		}
	case isDjangoView(prop.Module):
		if p := djangoViewRequestParam(f); isDjangoRequestParam(p) {
			updateSymbolWithDjangoRequest(ctx, p.Symbol, prop.Importer)
		}
	}
}

func doCallHeuristics(ctx fastnodectx.Context, callee pythontype.Value, args pythontype.Args, prop *propagator) []pythontype.Value {
	ctx.CheckAbort()

	if isDjangoGetModelCall(callee) {
		return doDjangoGetModelHeuristic(ctx, args, prop)
	}
	return nil
}
