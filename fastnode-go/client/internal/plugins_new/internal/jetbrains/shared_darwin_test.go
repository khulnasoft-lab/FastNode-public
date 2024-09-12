package jetbrains

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/plugins_new/editor"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/plugins_new/process"
	"github.com/stretchr/testify/require"
)

func newIntelliJTestManager(baseDir string, toolboxDir string, process *process.MockManager, betaChannel bool) (editor.Plugin, error) {
	manager, err := findManagerByIDTest(intellijID, process, betaChannel)
	if err != nil {
		return nil, err
	}

	manager.userHome = baseDir
	manager.toolboxDir = toolboxDir
	return manager, nil
}

func newPyCharmTestManager(baseDir string, toolboxDir string, process *process.MockManager, betaChannel bool) (editor.Plugin, error) {
	manager, err := findManagerByIDTest(pycharmID, process, betaChannel)
	if err != nil {
		return nil, err
	}

	manager.userHome = baseDir
	manager.toolboxDir = toolboxDir
	return manager, nil
}

func newGoLandTestManager(baseDir string, toolboxDir string, process *process.MockManager, betaChannel bool) (editor.Plugin, error) {
	manager, err := findManagerByIDTest(golandID, process, betaChannel)
	if err != nil {
		return nil, err
	}

	manager.userHome = baseDir
	manager.toolboxDir = toolboxDir
	return manager, nil
}

func newWebStormTestManager(baseDir string, toolboxDir string, process *process.MockManager, betaChannel bool) (editor.Plugin, error) {
	manager, err := findManagerByIDTest(webstormID, process, betaChannel)
	if err != nil {
		return nil, err
	}

	manager.userHome = baseDir
	manager.toolboxDir = toolboxDir
	return manager, nil
}

// setupIDEInstallation writes build.txt at the common location of the current OS and returns the ide home path
func setupIDEInstallation(t *testing.T, buildID string, path ...string) string {
	// on mac the toolbox stores applications inside the XY-123.456.78/ folder as an .app directory
	ideDir := filepath.Join(append(path, "IDE.app")...)
	err := os.MkdirAll(ideDir, 0700)
	require.NoError(t, err)

	buildFilePath := buildFileLocation(ideDir)
	os.MkdirAll(filepath.Dir(buildFilePath), 0700)

	err = ioutil.WriteFile(buildFilePath, []byte(buildID), 0600)
	require.NoError(t, err)
	return ideDir
}
