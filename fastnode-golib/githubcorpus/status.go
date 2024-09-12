package githubcorpus

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

// Stats ...
var (
	Stats                 = status.NewSection("githubcorpus")
	GetContentSuccessRate = Stats.Ratio("GetContentSuccessRate")
	GetCommitSuccessRate  = Stats.Ratio("GetCommitSuccessRate")
	MergeCommitRatio      = Stats.Ratio("MergeCommitRatio")
	FullCommitSuccessRate = Stats.Ratio("FullCommitSuccessRate")
)
