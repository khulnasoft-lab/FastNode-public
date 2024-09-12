package python

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythontype"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

func sbChoice(sbs []symbolBundle) {
	for i := range sbs {
		if i == 0 {
			continue
		}
		if pythontype.MoreSpecific(sbs[i].ns.val, sbs[0].ns.val) {
			sbs[i], sbs[0] = sbs[0], sbs[i]
		}
	}
}

func (vb valueBundle) chooseOne(ctx fastnodectx.Context) valueBundle {
	// we can directly construct a valueBundle here rather than using newValueBundle,
	// as the chosen value must have already been translated when client code called
	// newValueBundle to create the input.
	// This allows us to avoid excess calls to pythontype.Translate, etc
	return valueBundle{
		val:         pythontype.MostSpecific(ctx, vb.val),
		indexBundle: vb.indexBundle,
	}
}

func (sb symbolBundle) chooseOne(ctx fastnodectx.Context) symbolBundle {
	return symbolBundle{
		valueBundle: sb.valueBundle.chooseOne(ctx),
		ns:          sb.ns,
		nsName:      sb.nsName,
		name:        sb.name,
	}
}
