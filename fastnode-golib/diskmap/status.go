package diskmap

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("diskmap")

	getCounter            = section.Counter("Get calls")
	getDiskCounter        = section.Counter("Get calls hitting disk")
	notFoundCounter       = section.Counter("Not found")
	valueSizeSample       = section.SampleByte("Value size")
	bytesReadPerGetSample = section.SampleByte("Bytes read per Get")
)
