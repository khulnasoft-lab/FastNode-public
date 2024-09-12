package client

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section       = status.NewSection("client/internal/fastnode")
	userCount     = section.Counter("Users")
	skipBreakdown = section.Breakdown("Skip event reasons")
)

func init() {
	userCount.Set(1)
	skipBreakdown.AddCategories("unsaved", "unsupported file", "not whitelisted", "file too large", "editor skip")
}
