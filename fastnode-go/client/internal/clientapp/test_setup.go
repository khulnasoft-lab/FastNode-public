package clientapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/client"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clienttelemetry"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/mockserver"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/updates"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource/keytypes"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/licensing"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/telemetry"
	"github.com/pkg/errors"
)

// StartEmptyTestEnvironment creates an empty project and starts the fastnoded client
func StartEmptyTestEnvironment(components ...component.Core) (*TestEnvironment, error) {
	p, err := NewTestEnvironment()
	if err != nil {
		return nil, err
	}

	err = p.StartPortNoDists(0, components...)
	return p, err
}

// StartDefaultTestEnvironment creates a project with preconfigured files and start the fastnoded client.
func StartDefaultTestEnvironment(loginUser bool, clientOpts *client.Options, components ...component.Core) (*TestEnvironment, error) {
	p, err := NewTestEnvironment()
	if err != nil {
		return nil, err
	}

	setupDefaultFiles(p)

	// if using default opts, then preload the builtin distribution to avoid flaky tests
	if clientOpts == nil {
		clientOpts = &client.Options{
			LocalOpts: fastnodelocal.Options{
				Dists: []keytypes.Distribution{
					keytypes.BuiltinDistribution3,
				},
			},
		}
	}

	err = p.StartPort(0, clientOpts, components...)
	if err != nil {
		return p, fmt.Errorf("client startup failed with error: %s", err.Error())
	}

	if loginUser {
		resp, err := p.FastnodedClient.SendLoginRequest("user@example.com", "secret", true)
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("SendLoginRequest failed with error: %s", err.Error())
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("SendLoginRequest failed with unexpected status: %d/%s", resp.StatusCode, resp.Status)
		}
	}

	return p, err
}

// returns a set of whitelisted, blacklisted and ignored files
func setupDefaultFiles(p *TestEnvironment) {
	// setup 5 sample files
	for i := 0; i < 5; i++ {
		f := filepath.Join(p.DataDirPath, fmt.Sprintf("file_%d.py", i))
		ioutil.WriteFile(f, []byte("import json"), 0600)
		p.Files = append(p.Files, f)
	}
}

// NewTestEnvironment creates a new, empty project with the given feature flags enabled in the underlying platform
func NewTestEnvironment() (*TestEnvironment, error) {
	// Use homedir because the default temp directories are filtered by Fastnode on some platforms
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	dataDir, err := ioutil.TempDir(usr.HomeDir, "fastnode-user-data")
	if err != nil {
		return nil, err
	}

	backend, err := mockserver.NewBackend(map[string]string{"user@example.com": "secret", "pro@example.com": "secret"})
	if err != nil {
		return nil, fmt.Errorf("NewBackend failed with error: %s", err.Error())
	}
	backend.SetUserPlan("pro@example.com", true)

	ctx, cancel := context.WithCancel(context.Background())
	return &TestEnvironment{
		ctx:         ctx,
		ctxCancel:   cancel,
		DataDirPath: dataDir,
		Backend:     backend,
	}, nil
}

// TestEnvironment is a test setup which uses the real fastnoded server
type TestEnvironment struct {
	ctx         context.Context
	ctxCancel   func()
	Server      *http.Server
	DataDirPath string
	Backend     *mockserver.MockBackendServer
	Fastnoded       *client.Client
	FastnodedClient *mockserver.FastnodedClient
	MockTracker *telemetry.MockClient
	MockUpdater *updates.MockManager
	Files       []string
}

// StartPortNoDists activates the project without any Python dists loaded by fastnodelocal
func (p *TestEnvironment) StartPortNoDists(port int, components ...component.Core) error {
	return p.StartPort(port, &client.Options{
		LicenseStore: licensing.NewStore(p.Backend.Authority().CreateValidator(), ""),
		LocalOpts: fastnodelocal.Options{
			Dists: []keytypes.Distribution{},
		},
	}, components...)
}

// StartPort activates the project and stars fastnoded on the given port. Use '0' to let fastnoded choose its own port.
func (p *TestEnvironment) StartPort(port int, customOpts *client.Options, components ...component.Core) error {
	var opts client.Options
	if customOpts != nil {
		opts = *customOpts
	}

	if opts.TestRootDir == "" {
		opts.TestRootDir = p.DataDirPath
	}
	if opts.LocalOpts.IndexedDir == "" {
		opts.LocalOpts.IndexedDir = opts.TestRootDir
	}
	if opts.LicenseStore == nil {
		opts.LicenseStore = licensing.NewStore(p.Backend.Authority().CreateValidator(), "")
	}

	fastnoded, server, mockTracker, err := StartTestClient(p.ctx, port, &opts, components...)
	if err != nil {
		return errors.Wrap(err, "client setup failed")
	}

	// run client main loop in background
	go func() {
		if err := fastnoded.Connect(p.Backend.URL.String()); err != nil {
			log.Fatalf("connect failed: %s", err.Error())
		}
	}()

	fastnodedClient := mockserver.NewFastnodedClient(fastnoded.URL)

	p.Fastnoded = fastnoded
	p.FastnodedClient = fastnodedClient
	p.Server = server
	p.MockTracker = mockTracker

	switch m := p.Fastnoded.Updater.(type) {
	case *updates.MockManager:
		p.MockUpdater = m
	default:
		return errors.New("mock updater not accessible")
	}

	// Wait for client to be initialized before returning
	return p.WaitForReady(10 * time.Second)
}

// WaitForReady waits for the client in the test environment to initialize
func (p *TestEnvironment) WaitForReady(timeout time.Duration) error {
	// Try to wait for a bit so that the client has time to start
	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if p.Fastnoded.TestReady() {
				return nil
			}
		}
	}
}

// WaitForNotReady waits for the client in the test environment to disconnect
func (p *TestEnvironment) WaitForNotReady(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if !p.Fastnoded.TestReady() {
				return nil
			}
		}
	}
}

// Close releases resources
func (p *TestEnvironment) Close() {
	p.ctxCancel()

	// shutdown fastnoded's HTTP server
	if p.Server != nil {
		if err := p.Server.Close(); err != nil {
			log.Printf("error shuttding down fastnoded HTTP server: %s", err.Error())
		}
	}

	// disconnect fastnoded from backend
	if p.Fastnoded != nil {
		p.Fastnoded.Shutdown()
	}

	// mock backend
	if p.Backend != nil {
		p.Backend.Close()
	}

	clienttelemetry.Close()

	os.RemoveAll(p.DataDirPath)
}

// TestFlush calls TestFlush on the component manager of fastnoded
func (p *TestEnvironment) TestFlush(ctx context.Context) {
	p.Fastnoded.TestComponentManager().TestFlush(ctx)
}

// SetOffline calls SetOffline on the network manager of fastnoded
func (p *TestEnvironment) SetOffline(offline bool) {
	p.Fastnoded.Network.SetOffline(true)
}
