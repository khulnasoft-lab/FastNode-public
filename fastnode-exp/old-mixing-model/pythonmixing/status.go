package pythonmixing

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("lang/python/pythonmixing")

	newFeaturesDuration = section.SampleDuration("NewFeatures")

	modelInferDuration = section.SampleDuration("Infer")
)
