package main

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/rollbar"
import rollbarLib "github.com/rollbar/rollbar-go"

func init() {
	rollbarLib.SetLogger(&rollbarLib.SilentClientLogger{})
	rollbar.SetLogDisabled(true)

	rollbar.SetToken("XXXXXXX")
	rollbar.SetEnvironment("production")
	rollbar.SetClientVersion(version)
}
