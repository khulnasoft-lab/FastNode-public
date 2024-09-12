package pythonmodels

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonexpr"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonkeyword"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonmodels/callprob"
)

// Mock returns empty Models.
func Mock() *Models {
	return &Models{
		Keyword:      &pythonkeyword.Model{},
		Expr:         &pythonexpr.ModelShard{},
		FullCallProb: &callprob.Model{},
	}
}
