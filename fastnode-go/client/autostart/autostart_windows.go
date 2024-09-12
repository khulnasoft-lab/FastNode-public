// +build !standalone

package autostart

import (
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/reg"
)

func setEnabled(enabled bool) error {
	if enabled {
		return reg.UpdateHKCURun()
	}
	return reg.RemoveHKCURun()
}
