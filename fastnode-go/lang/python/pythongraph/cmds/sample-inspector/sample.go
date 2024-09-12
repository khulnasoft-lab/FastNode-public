package main

import (
	"bytes"
	"fmt"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonast"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythongraph"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythongraph/rendered"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonparser"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

type saver struct {
	Saved []pythongraph.SavedBundle
}

func (s *saver) Save(sb pythongraph.SavedBundle) {
	s.Saved = append(s.Saved, sb)
}

type renderedSample struct {
	Name  string
	Graph rendered.GraphWithExtras
}

func renderSavedSamples(saved ...pythongraph.SavedBundle) ([]*renderedSample, error) {
	var rs []*renderedSample
	for i, s := range saved {
		if pythonast.IsNil(s.AST) {
			s.AST, _ = pythonparser.Parse(fastnodectx.Background(), []byte(s.Buffer), parseOpts)
		}

		var ast string
		if !pythonast.IsNil(s.AST) {
			var b bytes.Buffer
			pythonast.Print(s.AST, &b, "\t")
			ast = b.String()
		}

		graph, err := rendered.NewGraphWithExtras(s, ast)
		if err != nil {
			return nil, err
		}

		rs = append(rs, &renderedSample{
			Name:  fmt.Sprintf("%d-%s", i, s.Label),
			Graph: graph,
		})
	}

	return rs, nil
}
