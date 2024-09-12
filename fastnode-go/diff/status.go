package diff

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("diff")

	resendRatio = section.Ratio("Requesting client resend text (buffer mismatch)")
)
