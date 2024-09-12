package localtraining

import (
	"io/ioutil"
	"log"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/errors"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0/bpe"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/rollbar"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/workerpool"
)

// ExtractNewVocabEntries ...
func ExtractNewVocabEntries(kctx fastnodectx.Context, files []string, language lexicalv0.LangGroup, builder *bpe.Builder, numGo, iters int) error {
	langLexer, err := lexicalv0.NewLexer(language.Lexer)
	if err != nil {
		return errors.Wrapf(err, "unable to get lexer")
	}

	jobs := make([]workerpool.Job, 0, len(files))
	for _, f := range files {
		fClose := f
		jobs = append(jobs, func() error {
			// TODO: make panic handling part of the workerpool?
			defer func() {
				if ex := recover(); ex != nil {
					rollbar.PanicRecovery(ex)
				}
			}()

			contents, err := ioutil.ReadFile(fClose)
			if err != nil {
				log.Printf("error reading file '%s': %v", fClose, err)
				return nil
			}

			tokens, err := langLexer.Lex(contents)
			if err != nil {
				log.Printf("error lexing file '%s': %v", fClose, err)
				return nil
			}

			var toks []string
			for _, tok := range tokens {
				if parts, ok := langLexer.ShouldBPEEncode(tok); ok {
					toks = append(toks, parts...)
				}
			}

			builder.Add(toks)
			return nil
		})
	}

	// TODO: use fastnodectx?
	pool := workerpool.NewWithCtx(kctx.Context(), numGo)
	defer pool.Stop()

	pool.Add(jobs)

	if err := pool.Wait(); err != nil {
		return errors.Wrapf(err, "pool error")
	}

	// TODO: add fastnodectx
	err = builder.Merge(bpe.MergeOptions{
		Iterations:  iters,
		Logging:     true,
		Concurrency: numGo,
	})
	if err != nil {
		return errors.Wrapf(err, "unable to build vocab")
	}

	return nil
}
