package websandbox

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythoncomplete/driver"
)

// Options contains process-wide settings and objects
type Options struct {
	Services            *python.Services
	IDCCCompleteOptions driver.Options
	SandboxRecordMode   bool
}
