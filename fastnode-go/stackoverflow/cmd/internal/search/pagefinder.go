package search

import "github.com/khulnasoft-lab/fastnode/fastnode-go/stackoverflow"

// PageFinder is interface for retrieving so page based on id.
type PageFinder interface {
	Find(id int64) (*stackoverflow.StackOverflowPage, error)
}
