package main

import (
	"log"
	"os"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
)

func main() {
	logging := len(os.Args) >= 2 && os.Args[1] == "true"

	println("Starting server on port 1234...\n")
	println("Configure fastnoded with: export FASTNODED_PYTHON_REMOTE=\"127.0.0.1:1234\"\n")
	_, _, err := pythonresource.StartServerDefaultOpts(":1234", logging)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
