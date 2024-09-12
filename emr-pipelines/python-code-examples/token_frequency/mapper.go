package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythoncode"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/awsutil"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/text"
)

func main() {
	r := awsutil.NewEMRIterator(os.Stdin)
	w := awsutil.NewEMRWriter(os.Stdout)
	defer w.Close()

	one, _ := json.Marshal(1)

	for r.Next() {
		var snippet pythoncode.Snippet
		err := json.Unmarshal(r.Value(), &snippet)
		if err != nil {
			log.Fatal(err)
		}
		for _, t := range text.TokenizeWithoutCamelPhrases(text.RemovePunctuations(snippet.Code)) {
			w.Emit(t, one)
		}
	}
	if err := r.Err(); err != nil {
		log.Fatalln("error reading stdin:", err)
	}
}
