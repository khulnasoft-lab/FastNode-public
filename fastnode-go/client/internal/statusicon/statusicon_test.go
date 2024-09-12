package statusicon

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
)

func Test_Component(t *testing.T) {
	m := NewManager(nil)
	component.TestImplements(t, m, component.Implements{
		Initializer: true,
		Settings:    true,
		UserAuth:    true,
	})
}
