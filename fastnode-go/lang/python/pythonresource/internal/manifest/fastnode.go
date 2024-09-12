//go:generate go-bindata -pkg manifest -prefix ../.. ../../manifest.json

package manifest

import (
	"bytes"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource/keytypes"
)

// FastnodeManifest is the global read-only manifest for all Fastnode resources
var FastnodeManifest Manifest

type emptyName struct {
	*bytes.Reader
}

func (r emptyName) Name() string { return "" }

func init() {
	m, err := New(emptyName{bytes.NewReader(MustAsset("manifest.json"))})
	if err != nil {
		panic(err)
	}

	FastnodeManifest = m

	if _, ok := FastnodeManifest[keytypes.BuiltinDistribution3]; !ok {
		panic("Python 3 builtin distribution not found in FastnodeManifest")
	}
}
