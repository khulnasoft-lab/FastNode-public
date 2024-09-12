package fastnodelocal

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section   = status.NewSection("client/internal/fastnodelocal")
	userCount = section.Counter("Users")
)

func init() {
	userCount.Set(1)
}
