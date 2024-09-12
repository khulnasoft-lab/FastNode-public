package pythonkeyword

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("lang/python/pythonkeyword")

	newFeaturesDuration = section.SampleDuration("NewFeatures")

	modelInferDuration  = section.SampleDuration("Model.Infer")
	modelIsKeywordRatio = section.Ratio("Model.Infer IsKeyword")
)
