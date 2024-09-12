package lexicalproviders

import (
	"sync"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0/golang"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0/javascript"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0/python"
)

var (
	configLock sync.Mutex
	pyConfig   python.Config
	jsConfig   javascript.Config
	goConfig   golang.Config
)

func init() {
	pyConfig = python.DefaultPrettifyConfig
	jsConfig = javascript.DefaultPrettifyConfig
	goConfig = golang.DefaultPrettifyConfig
}
