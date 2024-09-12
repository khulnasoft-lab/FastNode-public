package clientapp

import (
	"log"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/platform"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/platform/machine"
)

var (
	machineIDMessage = `Fastnode was unable to read its Machine ID from the registry. Please try reinstalling Fastnode. Fastnode will now exit.`
)

// Alert shows an alert UI for the given error.
func Alert(err error) {
	log.Println("alert:", err)
	switch err {
	case ErrPortInUse:
		// don't show a message
		// https://github.com/khulnasoft-lab/issue-tracker/issues/197
	case machine.ErrNoMachineID:
		platform.ShowAlert(machineIDMessage)
	default:
		platform.ShowAlert(err.Error())
	}
}
