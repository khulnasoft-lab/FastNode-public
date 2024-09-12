package desktoplogin

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section                 = status.NewSection("client/internal/desktoplogin")
	usedCounterDistribution = section.CounterDistribution("Desktop login used")
)
