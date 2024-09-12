package popularsignatures

import (
	"reflect"
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/editorapi"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/reflection"
)

func TestCast_Entity(t *testing.T) {
	if !reflection.StructurallyEqual(reflect.TypeOf([]*editorapi.Signature{}), reflect.TypeOf(Entity{})) {
		t.Logf("Entity type not structurally equal to []*editorapi.Signature")
		t.Fail()
	}
}
