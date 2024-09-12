module github.com/khulnasoft-lab/fastnode/windows

go 1.15

replace github.com/khulnasoft-lab/fastnode => ../

replace github.com/khulnasoft-lab/fastnode/fastnode-go/client/datadeps => ../fastnode-go/client/datadeps

require github.com/khulnasoft-lab/fastnode v0.0.0-00010101000000-000000000000 // indirect
