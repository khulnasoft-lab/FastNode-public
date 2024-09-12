package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythoncode"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/awsutil"
)

func main() {
	r := awsutil.NewEMRIterator(os.Stdin)
	w := awsutil.NewEMRWriter(os.Stdout)
	defer w.Close()

	for r.Next() {
		var usages []*pythoncode.ObjectUsage
		err := json.Unmarshal(r.Value(), &usages)
		if err != nil {
			log.Fatalln(err)
		}

		for _, usage := range usages {
			buf, err := json.Marshal(usage)
			if err != nil {
				log.Fatalln(err)
			}
			err = w.Emit(usage.Identifier, buf)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if err := r.Err(); err != nil {
		log.Fatalln(err)
	}
}
