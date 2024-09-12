package symgraph

import (
	"reflect"
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonimports"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/reflection"
)

func TestCast_DottedPath(t *testing.T) {
	if !reflection.StructurallyEqual(reflect.TypeOf(pythonimports.DottedPath{}), reflect.TypeOf(DottedPath{})) {
		t.Logf("DottedPath type not structurally equal to pythonimports.DottedPath")
		t.Fail()
	}
}
