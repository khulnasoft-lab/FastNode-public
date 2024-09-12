package codebase

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/localpath"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/recommend"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

type mockRecommender struct {
	files         []recommend.File
	shouldRebuild bool
}

func (r mockRecommender) Recommend(ctx fastnodectx.Context, request recommend.Request) ([]recommend.File, error) {
	return r.files, nil
}

func (r mockRecommender) RecommendBlocks(ctx fastnodectx.Context, request recommend.BlockRequest) ([]recommend.File, error) {
	return request.InspectFiles, nil
}

func (r mockRecommender) RankedFiles() ([]recommend.File, error) {
	return nil, nil
}

func (r mockRecommender) ShouldRebuild() (bool, error) {
	return r.shouldRebuild, nil
}

type mockIgnorer struct {
	shouldRebuild bool
}

func (i mockIgnorer) Ignore(pathname localpath.Absolute, isDir bool) bool {
	return false
}

func (i mockIgnorer) ShouldRebuild() (bool, error) {
	return i.shouldRebuild, nil
}
