package process

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

var attributes = &syscall.SysProcAttr{}
var bundleID = "com.fastnode.Fastnode"

// Name of Fastnode process
var Name = "Fastnode"

// Start attempts to start Fastnode.
func Start() error {
	loc, err := bundleLocation()
	if err != nil {
		return err
	}

	_, err = startProcess("open", nil, "-a", loc, "--args", "--plugin-launch")
	return err
}

func bundleLocation() (string, error) {
	out, err := exec.Command("mdfind", "kMDItemCFBundleIdentifier", "=", bundleID).Output()
	if err != nil {
		return "", err
	}

	var valid []string
	for _, x := range strings.Split(string(out), "\n") {
		if x != "" {
			valid = append(valid, x)
		}
	}

	if len(valid) < 1 {
		return "", errors.New("Couldn't find bundle location")
	}
	return valid[0], nil
}
