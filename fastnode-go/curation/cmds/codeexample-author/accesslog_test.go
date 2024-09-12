package main

import (
	"testing"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/curation"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoAccessLock(t *testing.T) {
	access := setupAccessManager()
	lock, err := access.currentAccessLock("lang", "pkg")
	assert.Nil(t, lock)
	assert.NoError(t, err)
}

func TestAcquireAccessLock(t *testing.T) {
	access := setupAccessManager()
	lock, err := access.acquireAccessLock("lang", "pkg", "test@khulnasoft.com")
	require.NotNil(t, lock)
	require.NoError(t, err)
	assert.Equal(t, "test@khulnasoft.com", lock.UserEmail)
}

func TestLockedOut(t *testing.T) {
	access := setupAccessManager()

	lockFoo, err := access.acquireAccessLock("lang", "pkg", "foo@khulnasoft.com")
	require.NotNil(t, lockFoo)
	require.NoError(t, err)

	lockBar, err := access.acquireAccessLock("lang", "pkg", "bar@khulnasoft.com")
	require.NotNil(t, lockBar)
	require.NoError(t, err)

	assert.Equal(t, lockFoo.UserEmail, lockBar.UserEmail)
}

func TestTakeOldAccessLock(t *testing.T) {
	access := setupAccessManager()
	access.timeout = 1

	lockFoo, err := access.acquireAccessLock("lang", "pkg", "foo@khulnasoft.com")
	require.NotNil(t, lockFoo)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	lockBar, err := access.acquireAccessLock("lang", "pkg", "bar@khulnasoft.com")
	require.NotNil(t, lockBar)
	require.NoError(t, err)
	assert.Equal(t, "bar@khulnasoft.com", lockBar.UserEmail)
}

func TestRenewAccessLock(t *testing.T) {
	access := setupAccessManager()
	access.timeout = 3

	lockFoo, err := access.acquireAccessLock("lang", "pkg", "foo@khulnasoft.com")
	require.NotNil(t, lockFoo)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	lockBar, err := access.acquireAccessLock("lang", "pkg", "bar@khulnasoft.com")
	require.NotNil(t, lockBar)
	require.NoError(t, err)
	assert.Equal(t, lockFoo.UserEmail, lockBar.UserEmail)

	lockFoo, err = access.acquireAccessLock("lang", "pkg", "foo@khulnasoft.com")
	require.NotNil(t, lockFoo)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	lockBar, err = access.acquireAccessLock("lang", "pkg", "bar@khulnasoft.com")
	require.NotNil(t, lockBar)
	require.NoError(t, err)
	assert.Equal(t, lockFoo.UserEmail, lockBar.UserEmail)
}

func setupAccessManager() *accessManager {
	db := curation.GormDB("sqlite3", ":memory:")
	access := newAccessManager(db)
	access.Migrate()
	return access
}
