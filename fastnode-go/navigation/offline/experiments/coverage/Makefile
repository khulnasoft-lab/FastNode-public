.PHONY: all clean

MAX_RECS_PER_FILE = -1
KHULNASOFT-LAB = $(GOPATH)/src/github.com/khulnasoft-lab/fastnode
GIT_CACHE = $(KHULNASOFT-LAB)/fastnode-go/navigation/offline/git-cache.json

all: results.md

results.md: commits-retrieved.csv text-retrieved.csv relevant.json analyze.py
	mypy --strict analyze.py
	python analyze.py \
		--max_tests=25 \
		--commits_retrieved_path=commits-retrieved.csv \
		--text_retrieved_path=text-retrieved.csv \
		--relevant_path=relevant.json \
		--histogram_path=histogram.png \
		--results_path=results.md

relevant.json: commits-retrieved.csv
	mypy --strict relevant.py
	python relevant.py \
		--repo_root=$(KHULNASOFT-LAB) \
		--git_cache=$(GIT_CACHE) \
		--retrieved_path=commits-retrieved.csv \
		--relevant_path=relevant.json

commits-retrieved.csv:
	go run ./cmds/retrieved \
		--usecommits \
		--maxrecsperfile=$(MAX_RECS_PER_FILE) \
		--reporoot=$(KHULNASOFT-LAB) \
		--gitcache=$(GIT_CACHE) \
		--retrievedpath=commits-retrieved.csv

text-retrieved.csv:
	go run ./cmds/retrieved \
		--maxrecsperfile=$(MAX_RECS_PER_FILE) \
		--reporoot=$(KHULNASOFT-LAB) \
		--gitcache=$(GIT_CACHE) \
		--retrievedpath=text-retrieved.csv

clean:
	-rm commits-retrieved.csv
	-rm text-retrieved.csv
	-rm relevant.json
