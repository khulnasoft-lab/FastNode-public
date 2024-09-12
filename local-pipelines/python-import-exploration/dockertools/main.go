package main

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/cmdline"
)

func main() {
	cmdline.MustDispatch(filesCmd, buildImageCmd, buildImagesCmd, deleteImageCmd, deleteImagesCmd)
}
