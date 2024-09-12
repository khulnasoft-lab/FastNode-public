// +build !windows

package platform

import "os"

// fastnodeRoot returns the directory containing fastnode configuration and session files.
func fastnodeRoot() string {
	return os.ExpandEnv("$HOME/.fastnode")
}
