package pythontype

import (
	"fmt"
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonimports"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

type mockValue struct {
	t    testing.TB
	kind Kind
	addr Address
}

// NewMockValue creates a new Value with the provided kind and addr for testing
func NewMockValue(t testing.TB, kind Kind, addr Address) Value {
	return mockValue{t, kind, addr}
}

// Kind implements Value
func (v mockValue) Kind() Kind {
	return v.kind
}

// Type implements Value
func (v mockValue) Type() Value {
	return NewMockValue(
		v.t,
		TypeKind,
		Address{
			File: "mockValueTypeAddrFile",
			Path: pythonimports.DottedPath{
				Hash: 12345,
				Parts: []string{
					"mock", "value", "type", "addr", "path",
				},
			},
		},
	)
}

// Address implements Value
func (v mockValue) Address() Address {
	return v.addr
}

// attr implements Value
func (v mockValue) attr(ctx fastnodectx.CallContext, name string) (AttrResult, error) {
	return AttrResult{}, fmt.Errorf("not implemented")
}

// Dir implements Value
func (v mockValue) Dir() ([]string, []Value) {
	return []string{}, []Value{}
}

// equal implements Value
func (v mockValue) equal(ctx fastnodectx.CallContext, other Value) bool {
	return v.Address().Equals(other.Address())
}

// Flatten implements Value
func (v mockValue) Flatten(flat *FlatValue, flattener *Flattener) {}

func (v mockValue) hash(ctx fastnodectx.CallContext) FlatID {
	return 0
}
