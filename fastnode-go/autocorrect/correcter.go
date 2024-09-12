package autocorrect

import (
	"sync"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/editorapi"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

var (
	correcters sync.Map
)

// Corrections is a set of corrections for a file.
type Corrections struct {
	NewBuffer string
}

// Correcter for corrections to a file.
type Correcter interface {
	Correct(ctx fastnodectx.Context, uid int64, mid string, req editorapi.AutocorrectRequest) (Corrections, error)
	Version() uint64
	ModelInfo(version uint64) (editorapi.AutocorrectModelInfoResponse, error)
}

// Register a correcter
func Register(correcter Correcter, l lang.Language) {
	correcters.Store(l, correcter)
}

func getCorrecter(ls string) (Correcter, bool) {
	var l lang.Language
	if ls == "python" {
		l = lang.Python
	}

	v, _ := correcters.Load(l)

	c, ok := v.(Correcter)

	return c, ok
}
