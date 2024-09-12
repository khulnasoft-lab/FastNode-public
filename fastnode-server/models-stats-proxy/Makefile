.PHONY: build push clean

build: tfserving-mux
	docker build -t khulnasoft/fastnode-server-models-stats-proxy .

push:
ifeq ($(TAG),)
	$(error Must set TAG to push)
endif
	docker tag khulnasoft-lab/fastnode-server-models-stats-proxy khulnasoft-lab/fastnode-server-models-stats-proxy:$(TAG)
	docker push khulnasoft/fastnode-server-models-stats-proxy:$(TAG)

clean:
	rm -rf tfserving-mux

tfserving-mux:
	env GOOS=linux go build github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0/tfserving/cmds/tfserving-mux/
