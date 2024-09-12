package main

import (
	"log"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/cmdline"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmdline.MustDispatch(countCmd, showCmd)
}
