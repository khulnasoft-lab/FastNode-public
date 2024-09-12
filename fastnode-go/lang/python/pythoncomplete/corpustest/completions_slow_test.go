// +build slow

package test

import (
	"testing"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/complete/corpustests"
)

func TestSlowCorpusTests(t *testing.T) {
	runFromCorpus(t, 30*time.Minute, corpustests.SlowState)
}
