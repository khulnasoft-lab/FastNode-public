.PHONY: build push clean

build: server
	docker build -t khulnasoft/fastnode-server-metadata .

push:
ifeq ($(TAG),)
	$(error Must set TAG to push)
endif
	docker tag khulnasoft-lab/fastnode-server-metadata khulnasoft-lab/fastnode-server-metadata:$(TAG)
	docker push khulnasoft/fastnode-server-metadata:$(TAG)

clean:
	rm -rf server

service/service_grpc.pb.go: service.proto
	protoc service.proto --go-grpc_out=.

server: service/service_grpc.pb.go server.go
	env GOOS=linux go build server.go
