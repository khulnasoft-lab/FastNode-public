package callprobcallmodel

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("lang/python/pythonmodels/callprob")

	newFeaturesDuration = section.SampleDuration("NewFeatures")

	modelInferDuration = section.SampleDuration("Infer")
)
