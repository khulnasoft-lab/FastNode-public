package javascript

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("lang/javascript (editor API)")

	tokensStatusCode = section.Breakdown("Tokens endpoint status codes")
	hoverStatusCode  = section.Breakdown("Hover endpoint status codes")
	calleeStatusCode = section.Breakdown("Callee endpoint status codes")
)
