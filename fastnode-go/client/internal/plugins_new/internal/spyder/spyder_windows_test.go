package spyder

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/plugins_new/editor"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/plugins_new/process"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Manager(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fastnode-spyder")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	iniData, err := ioutil.ReadFile(filepath.Join(prefix, "spyder.ini"))
	require.NoError(t, err)

	iniFilePath := filepath.Join(tempDir, ".spyder-py3", "config", "spyder.ini")
	err = os.MkdirAll(filepath.Dir(iniFilePath), 0700)
	require.NoError(t, err)

	condaDirPath := filepath.Join(tempDir, "anaconda3")
	condaPythonPath := filepath.Join(condaDirPath, "python.exe")
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Dir(condaPythonPath), 0700)
	require.NoError(t, err)
	err = ioutil.WriteFile(condaPythonPath, []byte{}, 0700)
	require.NoError(t, err)

	condaBatPath := filepath.Join(condaDirPath, "condabin", "conda.bat")
	err = os.MkdirAll(filepath.Dir(condaBatPath), 0700)
	require.NoError(t, err)
	err = ioutil.WriteFile(condaBatPath, []byte{}, 0700)
	require.NoError(t, err)

	err = ioutil.WriteFile(iniFilePath, iniData, 0600)
	require.NoError(t, err)

	p := process.MockManager{
		CustomDir: func() (string, error) {
			return tempDir, nil
		},
		RunResult: func(name string, arg ...string) ([]byte, error) {
			if name == condaBatPath {
				return []byte(` 
						[{
						"base_url": "https://conda.anaconda.org/spyder-ide",
						"build_number": 0,
						"build_string": "py37_0",
						"channel": "spyder-ide",
						"dist_name": "spyder-4.0.1-py37_0",
						"name": "spyder",
						"platform": "linux-64",
						"version": "4.0.1"
					  	}]`), nil
			}
			return nil, errors.Errorf("unexpected command " + name)
		},
		StartMenuData: func() []string {
			return []string{condaPythonPath}
		},
	}

	mgr, err := NewManager(&p)
	require.NoError(t, err)

	editors, err := mgr.DetectEditors(context.Background())
	require.NoError(t, err)
	require.Len(t, editors, 1)

	config, err := mgr.EditorConfig(context.Background(), editors[0])
	require.NoError(t, err)
	require.Empty(t, config.Compatibility)
	require.EqualValues(t, "4.0.0", config.RequiredVersion)
	require.EqualValues(t, "4.0.1", config.Version)
	require.EqualValues(t, iniFilePath, config.Path)

	require.EqualValues(t, ID, mgr.ID())
	require.EqualValues(t, name, mgr.Name())

	testInstallUninstallUpdate(t, mgr, iniFilePath)

	// activate again, apply suboptimal settings and test the HTTP requests
	err = setFastnodeEnabled(iniFilePath, true)
	require.NoError(t, err)
	err = setSpyderConfigValue(iniFilePath, "editor", "automatic_completions_after_chars", "3")
	require.NoError(t, err)

	optimalSettings, running, err := SettingsStatus(context.Background(), mgr)
	require.NoError(t, err)
	require.False(t, optimalSettings)
	require.False(t, running)

	err = ApplyOptimalSettings(context.Background(), mgr)
	require.NoError(t, err)

	optimalSettings, running, err = SettingsStatus(context.Background(), mgr)
	require.NoError(t, err)
	require.True(t, optimalSettings, "after applying new settings, settings have to be reported as optimal")
	require.False(t, running)
}

func Test_IsRunning(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fastnode-spyder")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// setup conda.bat to make the command pass the isCondaCmd check
	condaCmd := filepath.Join(tempDir, "condabin", "conda.bat")
	err = os.MkdirAll(filepath.Dir(condaCmd), 0700)
	require.NoError(t, err)

	err = ioutil.WriteFile(condaCmd, []byte{}, 0700)
	require.NoError(t, err)

	p := process.MockManager{
		ListData: func() (process.List, error) {
			return process.List{
				process.NewMockProcess("pythonw", filepath.Join(tempDir, "pythonw.exe"), []string{
					filepath.Join(tempDir, "pythonw.exe"),
					filepath.Join(tempDir, "scripts", "spyder-script.py"),
				}),
			}, nil
		},
	}

	mgr, err := NewManager(&p)
	require.NoError(t, err)

	running, err := mgr.DetectRunningEditors(context.Background())
	require.NoError(t, err)
	require.Len(t, running, 1)

	config := mgr.InstallConfig(context.Background())
	require.True(t, config.Running)
}

func Test_IsNotRunning(t *testing.T) {
	p := process.MockManager{
		ListData: func() (process.List, error) {
			return process.List{}, nil
		},
	}

	mgr, err := NewManager(&p)
	require.NoError(t, err)
	running, err := mgr.DetectRunningEditors(context.Background())
	require.NoError(t, err)
	require.Empty(t, running)

	config := mgr.InstallConfig(context.Background())
	require.False(t, config.Running)
}

// test that install, update and uninstall succeed
func testInstallUninstallUpdate(t *testing.T, mgr editor.Plugin, configFilePath string) {
	err := mgr.Install(context.Background(), configFilePath)
	require.NoError(t, err, "installing must succeed")
	assert.True(t, mgr.IsInstalled(context.Background(), configFilePath), "plugin must be installed after a successful call of Install")

	err = mgr.Update(context.Background(), configFilePath)
	require.NoErrorf(t, err, "updating must succeed")
	assert.True(t, mgr.IsInstalled(context.Background(), configFilePath), "plugin must still be installed after Update")

	err = mgr.Uninstall(context.Background(), configFilePath)
	require.NoErrorf(t, err, "uninstalling must succeed")
	assert.False(t, mgr.IsInstalled(context.Background(), configFilePath), "plugin must be uninstalled after Uninstall")
}
