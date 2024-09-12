package main

import (
	"encoding/gob"
	"log"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/stackoverflow"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/stackoverflow/cmd/internal/search"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fileutil"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/languagemodel"
)

const (
	awsSearchEndpoint = "https://search-stackoverflow-search-dev-0-aookytr4hofuhykxmab7nnygzq.us-west-2.cloudsearch.amazonaws.com"
	dataDirPathS3     = "s3://fastnode-data/stackoverflow/ranking/2015-10-29_11-38-56-AM/"
)

// fastnodeSearcher is a searcher that uses the fastnode search tools.
type fastnodeSearcher struct {
	index      search.Index
	pageFinder search.PageFinder
	ranker     *search.Ranker
	ld         *search.LanguageDetector
	rf         *search.ResultFilter
}

func newFastnodeSearcher(pf search.PageFinder) fastnodeSearcher {
	index, err := search.NewCloudSearchIndex(awsSearchEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	f, err := fileutil.NewCachedReader(fileutil.Join(dataDirPathS3, "tagData"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	var tcd search.TagClassData
	err = decoder.Decode(&tcd)
	if err != nil {
		log.Fatal(err)
	}
	rf := search.NewResultFilter(tcd)

	f, err = fileutil.NewCachedReader(fileutil.Join(dataDirPathS3, "ldscorer"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder = gob.NewDecoder(f)
	var scorer languagemodel.Scorer
	err = decoder.Decode(&scorer)
	if err != nil {
		log.Fatal(err)
	}
	ld := search.NewLanguageDetector(&scorer, tcd)

	fModel, err := fileutil.NewCachedReader(fileutil.Join(dataDirPathS3, "model.json"))
	if err != nil {
		log.Fatal(err)
	}

	fDoc, err := fileutil.NewCachedReader(fileutil.Join(dataDirPathS3, "docCounts"))
	if err != nil {
		log.Fatal(err)
	}

	ranker, err := search.NewRanker(fModel, fDoc)
	if err != nil {
		log.Fatal(err)
	}

	return fastnodeSearcher{
		index:      index,
		pageFinder: pf,
		ranker:     ranker,
		ld:         ld,
		rf:         rf,
	}
}

// Search satisfies the searcher interface
func (ks fastnodeSearcher) Search(query string, st stackoverflow.SearchType, maxResults int) ([]int64, error) {
	// TODO(juan): memory leak!
	// lang, found := s.ld.Detect(query)
	// if !found {
	// 	query = lang + " " + query
	// }
	ids, err := ks.index.Search(query, st, maxResults)
	if err != nil {
		return nil, err
	}
	// TODO(juan): memory leak!
	// var pages []*stackoverflow.StackOverflowPage
	// for _, id := range ids {
	// 	page, err := s.pageFinder.Find(id)
	// 	if err != nil {
	// 		log.Println("error for id:" + strconv.FormatInt(id, 10) + ", error msg: " + err.Error())
	// 		continue
	// 	}
	// 	pages = append(pages, page)
	// }
	// pages = s.rf.Filter(query, lang, pages)
	// err = s.ranker.Rank(query, pages)
	// if err != nil {
	// 	return nil, err
	// }
	// results := make([]stackoverflow.SearchResult, len(pages))
	// for i, page := range pages {
	// 	results[i] = stackoverflow.SearchResult{
	// 		ID: page.GetQuestion().GetPost().GetId(),
	// 	}
	// }
	return ids, nil
}
