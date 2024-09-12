GOPATH ?= ${HOME}/go
ifeq ($(origin GOPATH), undefined)
	GOPATH := ${HOME}/go
endif
export GOPATH
KHULNASOFT-LAB := github.com/khulnasoft-lab/fastnode/
FASTNODEPATH := ${GOPATH}/src/${KHULNASOFT-LAB}

analyze: asttrace.json
	${FASTNODEPATH}/fastnode-exp/telemetry-analysis/asttrace/analyze.py $<

asttrace.bin:
	go build -o $@ ${KHULNASOFT-LAB}/fastnode-exp/telemetry-analysis/asttrace

asttrace.json: | asttrace.bin
	./asttrace.bin >$@
