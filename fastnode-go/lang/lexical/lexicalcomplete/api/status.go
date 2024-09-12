package api

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	// Stats ...
	Stats = status.NewSection("api")
	// CompletionDuration ...
	CompletionDuration = Stats.SampleDuration("completion duration")
)
