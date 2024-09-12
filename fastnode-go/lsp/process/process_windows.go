package process

import (
	"log"
	"path/filepath"
	"syscall"

	"github.com/winlabs/gowin32"
	"github.com/winlabs/gowin32/wrappers"
)

var attributes = &syscall.SysProcAttr{HideWindow: true, CreationFlags: wrappers.CREATE_NO_WINDOW}

// Name of Fastnode process
var Name = "fastnoded.exe"

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
	programFiles, err := gowin32.GetKnownFolderPath(gowin32.KnownFolderProgramFiles)
	if err != nil {
		log.Println("error retrieving programFiles path", err)
		return "", err
	}
	// C:\Program Files\Fastnode
	p := filepath.Join(programFiles, "Fastnode", "fastnoded")
	return p, nil
}
