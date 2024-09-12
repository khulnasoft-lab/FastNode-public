package editorapi

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
)

func Test_Component(t *testing.T) {
	m := &Manager{}
	component.TestImplements(t, m, component.Implements{
		Initializer: true,
		Handlers:    true,
	})
}
