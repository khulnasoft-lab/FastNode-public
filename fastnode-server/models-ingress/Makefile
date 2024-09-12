.PHONY: build push

build:
	docker build -t khulnasoft/fastnode-server-models-ingress .

push:
ifeq ($(TAG),)
	$(error Must set TAG to push)
endif
	docker tag khulnasoft-lab/fastnode-server-models-ingress khulnasoft-lab/fastnode-server-models-ingress:$(TAG)
	docker push khulnasoft/fastnode-server-models-ingress:$(TAG)
