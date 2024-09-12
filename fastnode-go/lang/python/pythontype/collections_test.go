package pythontype

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/stretchr/testify/require"
)

func TestNamedTupleBoundMethods(t *testing.T) {
	typ := NewNamedTupleType(StrConstant("Blah"), []string{"field1", "field2"}).(NamedTupleType)
	val := NewNamedTupleInstance(typ, []Value{IntConstant(1), StrConstant("abc")})

	var count AttrResult
	err := fastnodectx.Background().WithCallLimit(1, func(ctx fastnodectx.CallContext) (err error) {
		count, err = val.attr(ctx, "count")
		return
	})
	require.NoError(t, err)
	MustHash(count.Value()) // just checking that it doesn't panic

	var index AttrResult
	err = fastnodectx.Background().WithCallLimit(1, func(ctx fastnodectx.CallContext) (err error) {
		index, err = val.attr(ctx, "index")
		return
	})
	require.NoError(t, err)
	MustHash(index.Value()) // just checking that it doesn't panic
}
