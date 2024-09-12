//go:generate go-bindata -pkg distidx -prefix ../.. ../../index.json

package distidx

import "bytes"

// FastnodeIndex is the global Fastnode distribution index
var FastnodeIndex Index

func init() {
	i, err := New(bytes.NewReader(MustAsset("index.json")))
	if err != nil {
		panic(err)
	}

	FastnodeIndex = i
}
