package index

import (
	"os"

	explocalfiles "github.com/khulnasoft-lab/fastnode/fastnode-exp/localcode-analysis/localfiles"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonbatch"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonlocal"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/localcode"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/localfiles"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/pipeline/sample"
)

type LocalIndex struct {
	Index  *pythonlocal.SymbolIndex
	Corpus sample.Corpus
	Files  []sample.FileInfo // the files that were used to build the index
}

func (LocalIndex) SampleTag() {}

func NewLocalIndex(c sample.Corpus, rm pythonresource.Manager) (LocalIndex, error) {
	fis, err := c.List()
	if err != nil {
		return LocalIndex{}, err
	}

	// don't build the index with site-packages files to improve performance
	var nonSP []sample.FileInfo
	for _, fi := range explocalfiles.CategorizeLocalFiles(fis) {
		if !fi.IsSitePackages {
			nonSP = append(nonSP, fi.FileInfo)
		}
	}

	if len(nonSP) == 0 {
		return LocalIndex{}, pipeline.NewErrorAsError("no non-site-package files available")
	}

	files := make([]*localfiles.File, 0, len(nonSP))
	for _, fi := range nonSP {
		files = append(files, &localfiles.File{
			Name:          fi.Name,
			HashedContent: fi.Name,
		})
	}

	fs := newCorpusFS(c, nonSP)

	bl := &pythonbatch.BuilderLoader{
		Graph:   rm,
		Options: pythonbatch.DefaultOptions,
	}

	res, err := bl.Build(fastnodectx.Background(), localcode.BuilderParams{
		Filename:   files[0].Name, // just choose a random files
		FileSystem: fs,
		Files:      files,
		FileGetter: fs,
	})
	if err != nil {
		return LocalIndex{}, err
	}

	return LocalIndex{
		Index:  res.LocalArtifact.(*pythonlocal.SymbolIndex),
		Corpus: c,
		Files:  nonSP,
	}, nil
}

type corpusFS struct {
	corpus sample.Corpus
	files  map[string]struct{}
}

func newCorpusFS(c sample.Corpus, fis []sample.FileInfo) corpusFS {
	files := make(map[string]struct{})
	for _, fi := range fis {
		files[fi.Name] = struct{}{}
	}
	return corpusFS{
		corpus: c,
		files:  files,
	}
}

func (c corpusFS) Stat(path string) (localcode.FileInfo, error) {
	if _, found := c.files[path]; found {
		return localcode.FileInfo{IsDir: false, Size: 100}, nil
	}
	return localcode.FileInfo{}, os.ErrNotExist
}

func (c corpusFS) ListRecursive(ctx fastnodectx.Context, path string) ([]*localfiles.File, error) {
	panic("not implemented")
}

func (c corpusFS) Get(filename string) ([]byte, error) {
	return c.corpus.Get(filename)
}
