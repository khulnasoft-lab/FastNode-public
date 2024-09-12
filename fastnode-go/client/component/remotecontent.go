package component

import "github.com/khulnasoft-lab/fastnode/fastnode-go/conversion/remotecontent"

// RemoteContentManager ...
type RemoteContentManager interface {
	RemoteContent() remotecontent.RemoteContent
}
