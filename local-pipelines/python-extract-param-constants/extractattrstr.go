package main

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonast"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythoncode"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonpipeline"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline"
)

// AttrConsts represents the attr consts indexed by func
type attrConsts pythoncode.ConstInfo

// SampleTag implements pipeline.Sample
func (s attrConsts) SampleTag() {}

func (s attrConsts) hitAttr(constant string) {
	s[constant]++
}

// ExtractAttrStr extracts the string constants for attribute
func extractAttrStr(s pipeline.Sample) attrConsts {
	rast := s.(pythonpipeline.Resolved).RAST
	ac := make(attrConsts)

	pythonast.Inspect(rast.Root, func(node pythonast.Node) bool {
		if attr, ok := node.(*pythonast.AttributeExpr); ok {
			if str, ok := attr.Value.(*pythonast.StringExpr); ok {
				ac.hitAttr(str.Literal())
			}
		}
		return true
	})
	return ac
}
