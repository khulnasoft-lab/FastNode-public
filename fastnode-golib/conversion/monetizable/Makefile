bindata.go:
	mkdir serve
	cp $(GOPATH)/src/github.com/khulnasoft-lab/fastnode/fastnode-python/analysis/conversion-model/params.json serve/params.json
	go-bindata -o bindata.go -pkg monetizable serve/
	rm -rf serve
