package permissions

import (
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0"
	"github.com/stretchr/testify/assert"
)

func TestLexicalAllLangsSupported(t *testing.T) {
	for _, l := range lexicalv0.WebGroup.Langs {
		for _, ext := range l.Extensions() {
			assert.Contains(t, supportMap, "."+ext)
		}
	}

	for _, l := range lexicalv0.JavaPlusPlusGroup.Langs {
		for _, ext := range l.Extensions() {
			assert.Contains(t, supportMap, "."+ext)
		}
	}

	for _, l := range lexicalv0.CStyleGroup.Langs {
		for _, ext := range l.Extensions() {
			assert.Contains(t, supportMap, "."+ext)
		}
	}
}
