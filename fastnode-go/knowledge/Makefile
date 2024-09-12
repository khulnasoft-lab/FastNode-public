.PHONY: app extract clean

ifeq ($(OS), Windows_NT)
KHULNASOFT-LAB = $(GOPATH)\src\github.com\khulnasoft-lab\khulnasoft-lab
PULLS_DIR = $(HOME)\fastnode-navigation
CLOSED_PULLS_PATH = $(HOME)\.fastnode\github\khulnasoft-lab\khulnasoft-lab\pulls.json
OPEN_PULLS_DIR = $(PULLS_DIR)\open
else
KHULNASOFT-LAB = $(GOPATH)/src/github.com/khulnasoft-lab/fastnode/
PULLS_DIR = $(HOME)/fastnode-navigation
CLOSED_PULLS_PATH = $(HOME)/.fastnode/github/khulnasoft-lab/khulnasoft-lab/pulls.json
OPEN_PULLS_DIR = $(PULLS_DIR)/open
endif

NUM_PULLS = 1000

app: server/bindata.go
	go run ./cmds/index/ --dirs $(KHULNASOFT-LAB)
	go run ./cmds/app/ --closedpullsdir=$(CLOSED_PULLS_PATH) --openpullsdir=$(OPEN_PULLS_DIR) --root=$(KHULNASOFT-LAB)

server/bindata.go: server/templates/* server/static/*
	go-bindata -pkg knowledge -o server/bindata.go server/templates/ server/static/

csindex:
	go run ./cmds/index/ --dirs $(KHULNASOFT-LAB)

extract:
	-mkdir -p $(PULLS_DIR)
	go run ./cmds/github-data/ --writedir=$(PULLS_DIR)
	@python pypulls/pulls.py $(GITHUB_AUTH_TOKEN) khulnasoft-lab/khulnasoft-lab

clean:
	go run ./cmds/index/
	rm -rf $(HOME)/.fastnode/github
	rm -rf $(PULLS_DIR)
