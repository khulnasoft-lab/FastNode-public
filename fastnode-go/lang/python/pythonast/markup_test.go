package pythonast_test

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonast"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonparser"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkup(t *testing.T) {
	src := `
import xyz
abc.xyz = foo()
`

	expected := `
import <NAME>xyz</NAME>
<ATTR><NAME>abc</NAME>.xyz</ATTR> = <NAME>foo</NAME>()
`

	var opts pythonparser.Options
	ast, err := pythonparser.Parse(fastnodectx.Background(), []byte(src), opts)
	require.NoError(t, err)

	out := pythonast.Markup([]byte(src), ast, func(n pythonast.Node) (begin, end string) {
		switch n.(type) {
		case *pythonast.NameExpr:
			return "<NAME>", "</NAME>"
		case *pythonast.AttributeExpr:
			return "<ATTR>", "</ATTR>"
		}
		return "", ""
	})

	assert.EqualValues(t, expected, out)
}
