package pythonstatic

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythontype"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

// walkSubclasses calls the walk function for each subclass of the given class
func walkSubclasses(ctx fastnodectx.Context, c *pythontype.SourceClass, f func(*pythontype.SourceClass) bool) {
	ctx.CheckAbort()

	if f(c) {
		for _, subclass := range c.Subclasses {
			walkSubclasses(ctx, subclass, f)
		}
	}
}
