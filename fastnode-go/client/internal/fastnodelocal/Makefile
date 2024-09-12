distranking-bindata:
	go build github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/cmds/distribution-ranking
	./distribution-ranking -output=distscores.json
	rm distribution-ranking
	go generate github.com/khulnasoft-lab/fastnode/fastnode-go/client/distranking
	rm distscores.json
