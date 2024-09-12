package main

import (
	"path/filepath"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/envutil"
)

const defaultMachine = "default"

var defaultDockerCerts = filepath.Join(envutil.MustGetenv("HOME"), ".docker/machine/machines/default")
