package autocorrect

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/status"
)

var (
	section = status.NewSection("/fastnode/autocorrect")

	responseCodes = section.Breakdown("Autocorrect response codes")
)
