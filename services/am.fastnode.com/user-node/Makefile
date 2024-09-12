.PHONY: clean

KHULNASOFT-LAB=$${PWD%/khulnasoft-lab/**}/khulnasoft-lab
CWD_RELATIVE=$${PWD\#/**/khulnasoft-lab}
GO_IMAGE=golang:1.15.3-buster

build/user-node: clean build main.go
	docker run --rm -e "GOPRIVATE=github.com/khulnasoft-lab/*" \
		-v $(KHULNASOFT-LAB):/go/src/github.com/khulnasoft-lab/fastnode \
		-v $(PWD)/build:/build \
		-w /go/src/github.com/khulnasoft-lab/$(CWD_RELATIVE) \
		$(GO_IMAGE) go build -o /build/user-node .

build:
	mkdir -p build

clean:
	rm -rf build
