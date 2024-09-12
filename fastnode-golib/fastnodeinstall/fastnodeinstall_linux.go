package fastnodeinstall

import (
	"strings"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/exec"
)

// IsSystemdTimerEnabled returns true if updates are automatically downloaded and applied by a system service.
// It checks if Fastnode's systemd user service is active.
func IsSystemdTimerEnabled() (bool, error) {
	cmd := exec.Command("systemctl", "--user", "show", "fastnode-updater.timer", "--property", "ActiveState")
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	return isActiveOutput(outBytes), nil
}

func isActiveOutput(outBytes []byte) bool {
	// ActiveState=active indicates an active service
	return "ActiveState=active" == strings.TrimSpace(string(outBytes))
}
