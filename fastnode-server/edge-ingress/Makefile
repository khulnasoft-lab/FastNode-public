.PHONY: build push
TAG :=

build:
	docker build -t khulnasoft/fastnode-server-edge-ingress .

ssl/edge-ingress.crt:
ssl/edge-ingress.key:
	openssl req  -nodes -new -x509  -keyout ssl/edge-ingress.key -out ssl/edge-ingress.crt

push:
ifeq ($(TAG),)
	$(error Must set TAG to push)
endif
	docker tag  khulnasoft-lab/fastnode-server-edge-ingress khulnasoft-lab/fastnode-server-edge-ingress:$(TAG)
	docker push khulnasoft/fastnode-server-edge-ingress:$(TAG)
