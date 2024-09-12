package options

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonparser"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonscanner"
)

// Parse options to use across binaries.
var Parse = pythonparser.Options{
	Approximate: false,
	ErrorMode:   pythonparser.FailFast,
}

// Lex options to use.
var Lex = pythonscanner.Options{
	ScanComments: false,
	ScanNewLines: true,
}
