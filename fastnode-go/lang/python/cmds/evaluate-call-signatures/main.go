//go:generate bash -c "go-bindata $BINDATAFLAGS -o bindata.go templates/..."

package main

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/cmdline"
)

func main() {
	cmdline.MustDispatch(diffCmd, coverageCmd, viewDiffCmd)
}
