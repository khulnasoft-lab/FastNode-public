package main

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/cmdline"
)

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cmdline.MustDispatch(wordCountCmd, vocabGenCmd)
}
