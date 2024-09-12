#!/usr/bin/env make

# This variable is used to set the environment variable UIDEBUG
# in sidebar-bindata. It is initialized with the value of UIDEBUG
# from the environment but is explicitly set to 1 when
# usernode-debug-ui is run.
UI_DEBUG := $(UIDEBUG)

# Allow for a test backend to be passed into the makefile, but make sure that we have
# a sensible default too
REACT_APP_TEST_BACKEND := $(if $(REACT_APP_TEST_BACKEND),$(REACT_APP_TEST_BACKEND),$("https://staging.khulnasoft.com"))

GITCOMMIT := $(shell git rev-parse HEAD)

.PHONY: run-standalone

default: install-standalone

#####################################
#  Go build and verification tools  #
#####################################

install-ci-deps:
	# This target contains a minimal set of tools needed by CI.
	# Do not add things here lightly!
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get github.com/jteeuwen/go-bindata/...
	go get gotest.tools/gotestsum

install-deps: install-ci-deps
	# Protocol buffers
	go get github.com/golang/protobuf/proto
	go get github.com/golang/protobuf/protoc-gen-go

	# Install some utilities
	go install github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/printjson
	go install github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/importchanged

datadeps-bindata:
	go install github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/cmds/datadeps-bindata

build-datadeps:
	./scripts/build_datadeps.sh

generate:
	go generate ./...

test:
	# Run gotestsum with codecov reports for fastnode-go and fastnode-golib
	gotestsum --raw-command scripts/go_test_coverage ./fastnode-go/... ./fastnode-golib/...
	# Run gotestsum for checking build & test for local-pipelines (not part of codecov)
	gotestsum ./local-pipelines/...
    # Only run the data race checker on goroutine-heavy packages
	go test -race \
		./fastnode-go/sandbox \
		./fastnode-go/client/internal/client \
		./fastnode-go/client/internal/clientapp \
		./fastnode-go/health/cmds/healthd \
		./fastnode-go/core \
		./fastnode-go/lang/python/pythonlocal

# Linux only, run tests with libtcmalloc overriding malloc, free, ...
test-tcmalloc:
	LD_PRELOAD="${PWD}/linux/tcmalloc/libtcmalloc_minimal_debug.so" ${MAKE} test

build:
	go build -v ./fastnode-go/... ./fastnode-golib/... ./local-pipelines/... ./fastnode-answers/...

vet:
	# Run go-vet on all directories
	go vet ./fastnode-go/... ./fastnode-golib/... ./local-pipelines/... ./fastnode-answers/...

lint:
	true ./scripts/custom_lint.sh
	# Run golint only on files that are not auto-generated
	find fastnode-go fastnode-golib local-pipelines fastnode-answers -name "*.go" | grep -v ".pb.go" | grep -v "bindata.go" | grep -v "stackoverflow-xml.go" | grep -v "lsp/types/protocol.go" | xargs -I file golint file > /tmp/golint.test 2>&1
	cat /tmp/golint.test
	! test -s /tmp/golint.test

fmt:
	find fastnode-go fastnode-golib local-pipelines fastnode-answers -name "*.go" | grep -v "bindata.go" | grep -v ".*.pb.go" | grep -v "/corpus/go/.*.go" | xargs -I file goimports -l=true file > /tmp/gofmt.test 2>/dev/null
	cat /tmp/gofmt.test
	! test -s /tmp/gofmt.test

check-client-fatal:
	true git grep 'log.Fatal' ./fastnode-go/client/internal/ ':(exclude)*_test.go' ':(exclude)*/cmds/*' > /tmp/fatal.test 2>&1
	cat /tmp/fatal.test
	! test -s /tmp/fatal.test

bin-check:
	! git status --porcelain --untracked-files=no | sed s/".* "// | xargs -I f file ../f | grep -E '(ELF|x86)'

verify: fmt lint vet bin-check build test

pull-frontend-docker:
	docker pull khulnasoft-lab/build-frontend

install-libtensorflow:
	sudo rm -f /usr/local/lib/libtensorflow* || true
	curl -L "https://s3-us-west-1.amazonaws.com/fastnode-data/tensorflow/libtensorflow-cpu-`go env GOOS`-x86_64-1.15.0.tar.gz" | sudo tar -C /usr/local -xz

install-libtensorflow-avx2:
	sudo rm -f /usr/local/lib/libtensorflow* || true
	curl -L "https://s3-us-west-1.amazonaws.com/fastnode-data/tensorflow/libtensorflow-cpu-`go env GOOS`-x86_64-avx2-1.15.0.tar.gz" | sudo tar -C /usr/local -xz


#######################################
#  Webapp assets/bindata generation   #
#######################################

# Ref for seemingly extravagant npm invocations: https://github.com/imagemin/pngquant-bin/issues/52#issuecomment-260247356
webapp-deps: pull-frontend-docker
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm rebuild --quiet; npm uninstall --quiet; npm install --quiet"

webapp-tests: webapp-deps
	# TODO(tarak): Use the right commands to run the tests here?
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm run build-test"
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app -t khulnasoft/build-frontend npm test

webapp-build: webapp-deps
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm run build"

webapp-build-dev: webapp-deps
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-e "REACT_APP_ENV=development"\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm run build"

webapp-build-staging: webapp-deps
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-e "REACT_APP_BACKEND=https://staging.khulnasoft.com" -e "REACT_APP_ENV=staging"\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm run build"

webapp-build-prod: webapp-deps
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-e "REACT_APP_BACKEND=https://alpha.khulnasoft.com" -e "REACT_APP_ENV=production"\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm run build"

webapp-build-testing: webapp-deps
	docker run --rm -v "$(PWD)":/khulnasoft-lab -w /khulnasoft-lab/web/app\
		-e "REACT_APP_BACKEND=$(REACT_APP_TEST_BACKEND)" -e "REACT_APP_ENV=development"\
		-t khulnasoft/build-frontend\
		/bin/bash -c "npm config set //registry.npmjs.org/:_authToken=$(NPM_TOKEN); npm run build"

#######################################
#  fastnoded.exe: windows                 #
#######################################

force:

fastnoded.exe: force
	go build -buildmode=exe \
		-ldflags "-H windowsgui -X github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp.gitCommit=$(GITCOMMIT)" \
		github.com/khulnasoft-lab/fastnode/fastnode-go/client/cmds/fastnoded

WINDOWS_BUILD_VERSION ?= "9.9.9.9"

FastnodeSetup.exe: fastnoded.exe fastnode-lsp.exe
	mv fastnoded.exe windows/
	mv fastnode-lsp.exe windows/
	mkdir -p windows/installer/current_build_bin/out
	cd windows/installer && ./nant.bat -D:prevPatchVersion="${WINDOWS_PATCH_BASE}" -D:buildnumstring="${WINDOWS_BUILD_VERSION}" build

FastnodeUpdateInfo.xml: FastnodeSetup.exe
	@cd windows/tools/fastnode_update_signer_cmd/bin/Debug && ./FastnodeUpdateSignerCmd.exe ${WINDOWS_PASS}

FastnodePatchUpdateInfo.xml: FastnodeSetup.exe
	@[[ -n "${WINDOWS_PATCH_BASE}" ]] && cd windows/tools/fastnode_patch_update_signer_cmd/bin/Debug && ./FastnodePatchUpdateSignerCmd.exe ${WINDOWS_PASS}

fastnode-lsp.exe: force
	go build \
		-ldflags "-H windowsgui" \
		github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/cmds/fastnode-lsp

fastnode-windows: FastnodeSetup.exe FastnodeUpdateInfo.xml FastnodePatchUpdateInfo.xml

#######################################

install-standalone:
	./scripts/standalone.sh install

run-standalone:
	./scripts/standalone.sh run

run-web-node:
	go run github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/web-node/
