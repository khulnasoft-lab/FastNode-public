package lexicalproviders

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/complete/data"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// OutputFunc is the callback for returning results from a Provider
type OutputFunc func(fastnodectx.Context, data.SelectedBuffer, MetaCompletion)

// Provider is a function that provides Completions for a SelectedBuffer by passing them to an OutputFunc
type Provider interface {
	Provide(fastnodectx.Context, Global, Inputs, OutputFunc) error
	Name() data.ProviderName
}
