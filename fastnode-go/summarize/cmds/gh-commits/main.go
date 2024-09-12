package main

import (
	"sort"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/summarize/data"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/cmdline"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline/aggregator"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline/source"
)

const (
	maxFileSize           = 1 << 18       // 256kb
	maxRepoSizeCompressed = 1 << 30       // 1 gb
	maxOutFileSize        = 4 * (1 << 30) // 4gb
)

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cmdline.MustDispatch(extractCmd, shuffleCmd, splitCmd)
}

func newGitCommitsWriter(out, name string, opts aggregator.WriterOpts) *aggregator.Writer {
	opts.Compress = true
	opts.MaxFileSizeBytes = maxOutFileSize
	if opts.FilePrefix == "" {
		opts.FilePrefix = "gh-commits"
	}

	return aggregator.NewJSONWriter(opts, name, out)
}

func newGitCommitsSource(dir, name string, opts source.DatasetOpts) pipeline.Source {
	fs, err := aggregator.ListDir(dir)
	fail(err)
	sort.Strings(fs)

	return source.NewDataset(opts, name, source.JSONProcessFn(data.GitCommit{}), fs...)
}
