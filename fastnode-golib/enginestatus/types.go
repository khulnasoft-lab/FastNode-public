package enginestatus

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/presentation"

// Response encapsulates a Fastnode Status response.
type Response struct {
	Status string               `json:"status"`
	Short  string               `json:"short"`
	Long   string               `json:"long"`
	Button *presentation.Button `json:"button,omitempty"`
}
