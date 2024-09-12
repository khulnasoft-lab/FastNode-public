package pythonkeyword

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonscanner"
	"github.com/stretchr/testify/assert"
)

func TestKeywordMappingExists(t *testing.T) {
	for _, tok := range pythonscanner.KeywordTokens {
		//We check that all keywords have a mapping
		assert.NotEqual(t, AllKeywords[tok].Cat, 0)
	}
}
