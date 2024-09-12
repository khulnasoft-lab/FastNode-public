package autocorrect

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("/autocorrect")

	breakdown = section.Breakdown("Outcomes")

	segmentResults = section.Breakdown("segment results")
)
