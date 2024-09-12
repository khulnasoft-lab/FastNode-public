package test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IndexingStatus(t *testing.T) {
	project, err := startFastnodeLocal()
	require.NoError(t, err)
	defer project.Close()

	client := project.FastnodedClient

	// use temp dir within the project.DataDirPath so the files are walked and watched
	dir, err := ioutil.TempDir(project.DataDirPath, "fastnode-temp")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "file.py")

	remote := NewEditorRemoteControl("atom", project, t)
	remote.OpenNewFile(file)

	status, err := client.FileStatus(file)
	require.NoError(t, err)
	assert.EqualValues(t, "indexing", status, "expected fallback value 'indexing' for file without events")

	remote.Input("")
	project.TestFlush(context.Background())

	// An index with one file should have status 'ready'.
	// we need to wait until the indexer is ready, it's not guaranteed to be ready after the first edit event
	// we have to send distinct events to make this work
	for i := 0; i < 10; i++ {
		remote.Input("import json\n")
		project.TestFlush(context.Background())
		time.Sleep(time.Second)
		if status, err = client.FileStatus(file); status != "indexing" {
			break
		}
	}

	status, err = client.FileStatus(file)
	assert.EqualValues(t, "ready", status, "expected value 'ready' for file with events after several retries")

	// add another file
	file = filepath.Join(dir, "file2.py")
	remote.OpenNewFile(file)
	remote.Input("")
	project.TestFlush(context.Background())

	// An index with > 1 file should have status 'ready'.
	// we need to wait until the indexer is ready, it's not guaranteed to be ready after the first edit event
	// we have to send distinct events to make this work
	for i := 0; i < 10; i++ {
		remote.Input("import json\n")
		project.TestFlush(context.Background())
		time.Sleep(time.Second)
		if status, err = client.FileStatus(file); status == "ready" {
			break
		}
	}

	status, err = client.FileStatus(file)
	assert.EqualValues(t, "ready", status, "expected value 'ready' for file with events after several retries")
}

func Test_FilteredStatus(t *testing.T) {
	project, err := startFastnodeLocal()
	require.NoError(t, err)
	defer project.Close()

	client := project.FastnodedClient

	dir, err := ioutil.TempDir("", "fastnode-temp")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	// The default temp dir should be filtered on linux and windows
	if runtime.GOOS == "darwin" {
		dir = filepath.Join(dir, "Library")
		err = os.Mkdir(dir, os.ModePerm)
		require.NoError(t, err)
	}
	file := filepath.Join(dir, "file.py")

	// status should be noIndex once events have been flushed (file not written to disk)
	flushEvents(t, project)
	status, err := client.FileStatus(file)
	require.NoError(t, err)
	assert.EqualValues(t, "noIndex", status)
}

func Test_LocalLibAutosearch(t *testing.T) {
	project, err := startFastnodeLocal()
	require.NoError(t, err)
	defer project.Close()

	client := project.FastnodedClient

	autosearchClient, err := NewClient(project)
	require.NoError(t, err)
	defer autosearchClient.Close()

	// make fastnode libraries directory
	fastnodeLibs := filepath.Join(project.Fastnoded.Platform.FastnodeRoot, "libraries")
	err = os.Mkdir(fastnodeLibs, os.ModePerm)
	require.NoError(t, err)

	// create our library directory within site-packages
	tmpDir, err := ioutil.TempDir(project.Fastnoded.Platform.FastnodeRoot, "tmpDir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	sitePackages := filepath.Join(tmpDir, "site-packages")
	err = os.Mkdir(sitePackages, os.ModePerm)
	require.NoError(t, err)

	libDir := filepath.Join(sitePackages, "testlib")
	err = os.Mkdir(libDir, os.ModePerm)
	require.NoError(t, err)

	// add symlink to fastnode lib directory
	envPath := filepath.Join(fastnodeLibs, "foo") // where our fastnode lib will live
	err = os.Symlink(sitePackages, envPath)
	require.NoError(t, err)

	// add library files
	message := []byte("import json\nimport os\n")
	err = ioutil.WriteFile(filepath.Join(libDir, "__init__.py"), message, os.ModePerm)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(libDir, "test1.py"), message, os.ModePerm)
	require.NoError(t, err)

	// use temp dir within the project.DataDirPath so the files are walked and watched
	dir, err := ioutil.TempDir(project.DataDirPath, "fastnode-temp")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "file.py")

	remote := NewEditorRemoteControl("atom", project, t)
	remote.OpenNewFile(file)

	remote.Input("import testlib\n")
	remote.Save()
	project.TestFlush(context.Background())

	// An index with one file should have status 'ready'.
	// we need to wait until the indexer is ready, it's not guaranteed to be ready after the first edit event
	// we have to send distinct events to make this work
	for i := 0; i < 10; i++ {
		remote.Input("import re\n")
		project.TestFlush(context.Background())
		time.Sleep(time.Second)
		if status, _ := client.FileStatus(file); status != "indexing" {
			break
		}
	}

	status, err := client.FileStatus(file)
	assert.EqualValues(t, "ready", status, "expected value 'ready' for file with events after several retries")

	remote.Input("from testlib import os")

	autosearchID, err := autosearchClient.ReceiveClientMessage()
	require.NoError(t, err)
	require.NotEmpty(t, autosearchID)
	require.True(t, strings.Contains(autosearchID, "site-packages:testlib;;os"), "autosearch ID does not contain expected id. value: %s", autosearchID)

	members, err := client.Members(autosearchID)
	require.NoError(t, err)
	require.NotEmpty(t, members.Members, "members expected for libDir module")
}

func flushEvents(t *testing.T, project *clientapp.TestEnvironment) {
	// wait for fs watcher to pickup any pending change
	time.Sleep(1 * time.Second)

	fastnodeLocal := findFastnodeLocalComponent(project)
	require.NotNil(t, fastnodeLocal)
	fastnodeLocal.TestFlush(context.Background())
}

func findFastnodeLocalComponent(e *clientapp.TestEnvironment) *fastnodelocal.Manager {
	for _, c := range e.Fastnoded.Components() {
		if kl, ok := c.(*fastnodelocal.Manager); ok {
			return kl
		}
	}
	return nil
}
