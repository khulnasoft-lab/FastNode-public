package parser

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("lang/javascript (parser)")

	parseDuration     = section.SampleDuration("Parse duration")
	tooManyStepsRatio = section.Ratio("Parser took to many steps")
)
