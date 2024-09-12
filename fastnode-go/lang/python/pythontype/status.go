package pythontype

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	translationSection = status.NewSection("pythontype/translate")

	translateGlobalSuccesRatio = translationSection.Ratio("Translate global success ratio")

	translateGlobalFailures = translationSection.Breakdown("Translate global failures")
)
