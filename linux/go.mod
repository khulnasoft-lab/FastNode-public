module github.com/khulnasoft-lab/fastnode/linux

go 1.15

replace github.com/khulnasoft-lab/fastnode/fastnode-go/client/datadeps => ../fastnode-go/client/datadeps

replace github.com/khulnasoft-lab/fastnode => ../

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/khulnasoft-lab/go-bsdiff/v2 v2.0.1 // indirect
	github.com/khulnasoft-lab/fastnode v0.0.0-00010101000000-000000000000
	github.com/klauspost/cpuid v1.3.1
	github.com/mitchellh/cli v1.1.2
	github.com/rollbar/rollbar-go v1.2.0
	github.com/shirou/gopsutil v2.20.2+incompatible
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20201201145000-ef89a241ccb3
)
