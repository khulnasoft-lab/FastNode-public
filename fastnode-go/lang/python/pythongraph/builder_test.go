package pythongraph

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonanalyzer"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonparser"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonscanner"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/stretchr/testify/require"
)

var erm pythonresource.Manager

func emptyRM(t testing.TB) pythonresource.Manager {
	if erm == nil {
		erm = pythonresource.MockManager(t, nil)
	}
	return erm
}

func requireBuilderOpts(t *testing.T, rm pythonresource.Manager, src string, addMissingNames bool) *graphBuilder {
	sb := []byte(src)
	words, err := pythonscanner.Lex(sb, scanOpts)
	require.NoError(t, err)

	ast, err := pythonparser.ParseWords(fastnodectx.Background(), sb, words, parseOpts)
	require.NoError(t, err)

	rast, err := pythonanalyzer.NewResolver(rm, pythonanalyzer.Options{
		Path: "/src.py",
	}).Resolve(ast)

	require.NoError(t, err)

	return newBuilder(fastnodectx.Background(), newAnalysis(rm, words, rast), addMissingNames, false)
}

func requireBuilder(t *testing.T, rm pythonresource.Manager, src string) *graphBuilder {
	return requireBuilderOpts(t, rm, src, true)
}
