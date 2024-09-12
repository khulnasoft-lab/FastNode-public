package spyder

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

var prefix string

func init() {
	if runtime.GOOS == "windows" {
		prefix = filepath.Join("test", "windows")
	} else {
		prefix = filepath.Join("test", "unix")
	}
}

func Test_ConfigUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "fastnode-spyder")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	assertConfigFileUpdate(t, true, filepath.Join(prefix, "spyder.ini"), filepath.Join(prefix, "spyder.ini.true"))
	assertConfigFileUpdate(t, false, filepath.Join(prefix, "spyder.ini"), filepath.Join(prefix, "spyder.ini.false"))

	assertConfigFileUpdate(t, true, filepath.Join(prefix, "spyder_no_fastnode_enabled.ini"), filepath.Join(prefix, "spyder_no_fastnode_enabled.ini.true"))
	assertConfigFileUpdate(t, false, filepath.Join(prefix, "spyder_no_fastnode_enabled.ini"), filepath.Join(prefix, "spyder_no_fastnode_enabled.ini.false"))

	assertConfigFileUpdate(t, true, filepath.Join(prefix, "spyder_no_fastnode_section.ini"), filepath.Join(prefix, "spyder_no_fastnode_section.ini.true"))
	assertConfigFileUpdate(t, false, filepath.Join(prefix, "spyder_no_fastnode_section.ini"), filepath.Join(prefix, "spyder_no_fastnode_section.ini.false"))

	assertConfigFileUpdate(t, true, filepath.Join(prefix, "spyder_fastnode_eof.ini"), filepath.Join(prefix, "spyder_fastnode_eof.ini.true"))
	assertConfigFileUpdate(t, false, filepath.Join(prefix, "spyder_fastnode_eof.ini"), filepath.Join(prefix, "spyder_fastnode_eof.ini.false"))
}

func Test_IsFastnodeEnabled(t *testing.T) {
	dir, err := ioutil.TempDir("", "fastnode-spyder")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	assertFastnodeEnabled(t, true, filepath.Join(prefix, "spyder.ini"))
	assertFastnodeEnabled(t, true, filepath.Join(prefix, "spyder.ini.true"))
	assertFastnodeEnabled(t, false, filepath.Join(prefix, "spyder.ini.false"))

	assertFastnodeEnabled(t, false, filepath.Join(prefix, "spyder_no_fastnode_section.ini"))

	assertFastnodeEnabled(t, false, filepath.Join(prefix, "spyder_no_fastnode_enabled.ini"))
}

func Test_ApplyOptimalSettings(t *testing.T) {
	dir, err := ioutil.TempDir("", "fastnode-spyder")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config := filepath.Join(dir, "spyder.ini")
	err = ioutil.WriteFile(config, []byte(""), 0600)
	require.NoError(t, err)

	mayApply := couldApplyOptimalSettings(config)
	require.False(t, mayApply, "optimized settings must not be applicable if fastnode and completions are disabled")

	// replicate Spyder's defaults
	err = setFastnodeEnabled(config, false)
	require.NoError(t, err)
	err = setSpyderConfigValue(config, "editor", "automatic_completions", "True")
	require.NoError(t, err)
	err = setSpyderConfigValue(config, "editor", "automatic_completions_after_chars", "3")
	require.NoError(t, err)
	err = setSpyderConfigValue(config, "editor", "automatic_completions_after_ms", "300")
	require.NoError(t, err)

	mayApply = couldApplyOptimalSettings(config)
	require.False(t, mayApply, "optimized settings must not applicable with Fastnode disabled")

	err = setFastnodeEnabled(config, true)
	require.NoError(t, err)

	mayApply = couldApplyOptimalSettings(config)
	require.True(t, mayApply, "optimized settings must be applicable if completions and fastnode and enabled")

	err = applyOptimalSettings(config)
	require.NoError(t, err)

	chars, err := getSpyderConfigValue(config, "editor", "automatic_completions_after_chars")
	require.NoError(t, err)
	require.EqualValues(t, "1", chars)

	delay, err := getSpyderConfigValue(config, "editor", "automatic_completions_after_ms")
	require.NoError(t, err)
	require.EqualValues(t, "100", delay)
}

func assertConfigFileUpdate(t *testing.T, enabledValue bool, configFilePath string, refConfigFilePath string) {
	// copy test file to temp dir and modify it
	iniData, err := ioutil.ReadFile(configFilePath)
	require.NoError(t, err)

	tempFile, err := ioutil.TempFile("", "spyder-ini")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	tempFilePath := tempFile.Name()
	_, err = tempFile.Write(iniData)
	require.NoError(t, err)
	_ = tempFile.Close()

	err = setFastnodeEnabled(tempFilePath, enabledValue)
	require.NoError(t, err)

	// compare update data with expected data
	tempFileData, err := ioutil.ReadFile(tempFilePath)
	require.NoError(t, err)
	refData, err := ioutil.ReadFile(refConfigFilePath)
	require.NoError(t, err)
	require.EqualValues(t, string(refData), string(tempFileData))
}

func assertFastnodeEnabled(t *testing.T, expected bool, configFilePath string) {
	isEnabled := isFastnodeEnabled(configFilePath)
	require.EqualValues(t, expected, isEnabled)
}
