package main

import (
	"github.com/alexflint/go-arg"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/codesearch"
)

func main() {
	args := struct {
		Dirs []string
	}{}
	arg.MustParse(&args)
	opts := codesearch.IndexOptions{
		ResetFlag: true,
	}
	codesearch.Index(opts, args.Dirs...)
}
