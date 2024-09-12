package process

import (
	"path/filepath"
	"syscall"
)

var attributes = &syscall.SysProcAttr{}

// Name of Fastnode process
var Name = "fastnoded"

// Start attempts to start Fastnode.
func Start() error {
	path, err := installPath()
	if err != nil {
		return err
	}
	_, err = startProcess(path, nil, "--plugin-launch")
	if err != nil {
		return err
	}
	return nil
}

// Use default location defined in https://help.khulnasoft.com/article/136-how-to-restart-fastnode
func installPath() (string, error) {
	hd, err := homeDir()
	if err != nil {
		return "", err
	}
	// ~/.local/share/fastnode/fastnoded
	p := filepath.Join(hd, ".local", "share", "fastnode", "fastnoded")
	return p, nil
}
