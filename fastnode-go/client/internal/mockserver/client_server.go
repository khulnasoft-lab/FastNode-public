package mockserver

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/config"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clienttelemetry"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/conversion/cohort"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/conversion/monetizable"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/permissions"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/metrics"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/settings"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/platform"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/community"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/community/account"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/userids"
)

// NewTestClientServer returns a new TestClientServer which accepts the given validUsers.
// Returns an error if the setup failed
func NewTestClientServer(validUsers map[string]string) (*TestClientServer, error) {
	return NewTestClientServerFeatures(validUsers, nil)
}

// NewTestClientServerFeatures returns a new TestClientServer which accepts the given validUsers.
// Returns an error if the setup failed
func NewTestClientServerFeatures(validUsers map[string]string, featureOverride map[string]bool) (*TestClientServer, error) {
	return NewTestClientServerRootFeatures("", validUsers, featureOverride)
}

// NewTestClientServerRootFeatures returns a new TestClientServer which accepts the given Fastnode root & validUsers.
// Returns an error if the setup failed
func NewTestClientServerRootFeatures(root string, validUsers map[string]string, featureOverride map[string]bool) (*TestClientServer, error) {
	fastnoded, err := NewFastnodedTestServer()
	if err != nil {
		return nil, err
	}

	backend, err := NewBackend(validUsers)
	if err != nil {
		return nil, err
	}

	fastnodedClient := NewFastnodedClient(fastnoded.URL)

	p, err := platform.NewTestPlatformFeatures(root, featureOverride)
	if err != nil {
		return nil, err
	}

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	basePath, err := ioutil.TempDir(usr.HomeDir, "fastnode-temp-files")
	if err != nil {
		return nil, err
	}

	t := &TestClientServer{
		Fastnoded:                 fastnoded,
		FastnodedClient:           fastnodedClient,
		Backend:               backend,
		Components:            component.NewTestManager(),
		Platform:              p,
		quitChan:              make(chan int, 1),
		ReadLoginLogoutEvents: true,
		Languages:             []lang.Language{lang.Python},
		BasePath:              basePath,
	}
	if root == "" {
		t.AddCleanupAction(func() {
			os.RemoveAll(t.Platform.FastnodeRoot)
		})
	}
	return t, nil
}

// TestClientServer bundles a mock backend server, a mock fastnoded server and a client to talk to the fastnoded server
type TestClientServer struct {
	Fastnoded       *FastnodedTestServer
	Backend     *MockBackendServer
	Components  *component.Manager
	AuthClient  component.AuthClient
	Permissions component.PermissionsManager
	Settings    component.SettingsManager
	Metrics     component.MetricsManager
	Platform    *platform.Platform
	Network     component.NetworkManager

	// languages supported by our fastnoded mock, defaults to Python
	Languages []lang.Language

	FastnodedClient *FastnodedClient

	quitChan chan int

	ReadLoginLogoutEvents bool

	cleanupActions []func()

	// guards fastnodeUser
	mu       sync.Mutex
	fastnodeUser *community.User

	BasePath string
}

// AddCleanupAction registers adds a new cleanup action which will be called by the Close() method
func (t *TestClientServer) AddCleanupAction(a func()) {
	t.cleanupActions = append(t.cleanupActions, a)
}

// Close releases resources used by the TestClientServer
func (t *TestClientServer) Close() {
	t.quitChan <- 0

	t.Components.Terminate()
	t.Fastnoded.Close()
	t.Backend.Close()

	for _, a := range t.cleanupActions {
		a()
	}

	os.RemoveAll(t.BasePath)
}

// SendAccountCreationRequest posts the email and password to the server. It optionally waits until the account creation has been processed by the mockserver
func (t *TestClientServer) SendAccountCreationRequest(email, password string, waitForCreation bool) (*http.Response, error) {
	return t.FastnodedClient.SendAccountCreationRequest(email, password, waitForCreation)
}

// SendLoginRequest posts the email and password to the server. It optionally waits until the login has been processed by the mockserver
func (t *TestClientServer) SendLoginRequest(email, password string, waitForLogin bool) (*http.Response, error) {
	return t.FastnodedClient.SendLoginRequest(email, password, waitForLogin)
}

// SendLogoutRequest is a helper method to logout
func (t *TestClientServer) SendLogoutRequest(waitForLogout bool) (*http.Response, error) {
	return t.FastnodedClient.SendLogoutRequest(waitForLogout)
}

// DoFastnodedGet is a helper method to send a HTTP GET request to fastnoded
func (t *TestClientServer) DoFastnodedGet(path string) (*http.Response, error) {
	return t.FastnodedClient.Get(path)
}

// DoFastnodedPost is a helper method to send a HTTP POST request to fastnoded
func (t *TestClientServer) DoFastnodedPost(path string, body io.Reader) (*http.Response, error) {
	return t.FastnodedClient.Post(path, body)
}

// DoFastnodedPut is a helper method to send a HTTP PUT request to fastnoded
func (t *TestClientServer) DoFastnodedPut(path string, body io.Reader) (*http.Response, error) {
	return t.FastnodedClient.Put(path, body)
}

// MockNetworkManager is a mock NetworkManager, who's needed to avoid a
// network->auth->mockserver->network import cycle
type MockNetworkManager struct {
	online      bool
	fastnodedOnline bool
}

// Name implements interface Core
func (m *MockNetworkManager) Name() string {
	return "network"
}

// SetOnline sets the network to online or offline based on the value of the bool
func (m *MockNetworkManager) SetOnline(val bool) {
	m.online = val
}

// SetOffline sets the network to online or offline based on the value of the bool
func (m *MockNetworkManager) SetOffline(val bool) {
	m.online = !val
}

// Online implements interface NetworkManager
func (m *MockNetworkManager) Online() bool {
	return m.online
}

// CheckOnline implements interface component.NetworkManager
func (m *MockNetworkManager) CheckOnline(ctx context.Context) bool {
	return m.online
}

// FastnodedOnline implements interface component.NetworkManager
func (m *MockNetworkManager) FastnodedOnline() bool {
	return m.fastnodedOnline
}

// FastnodedInitialized implements interface FastnodedEventer
func (m *MockNetworkManager) FastnodedInitialized() {
	m.fastnodedOnline = true
}

// FastnodedUninitialized implements interface FastnodedEventer
func (m *MockNetworkManager) FastnodedUninitialized() {
	m.fastnodedOnline = false
}

// NewMockNetworkManager returns a new MockNetworkManager
func NewMockNetworkManager() component.NetworkManager {
	return &MockNetworkManager{
		online: true,
	}
}

// SetupWithCustomAuthClient configures performs a default setup, but uses a custom auth client
// This is useful for tests which would otherwise create an import cycle, e.g. to test metrics which uses the TestClientServer
// which used metrics if there wasn't a way to pass a custom auth client
func (t *TestClientServer) SetupWithCustomAuthClient(authClient component.AuthClient, components ...component.Core) error {
	s := settings.NewTestManager()

	metrics := metrics.NewMockManager()
	permMgr := permissions.NewManager(t.Languages, nil)

	return t.SetupComponents(authClient, s, permMgr, metrics, components...)
}

// SetupComponents configures the components registered in the mocked Fastnoded http server
func (t *TestClientServer) SetupComponents(auth component.AuthClient, settings component.SettingsManager, permissions component.PermissionsManager, metrics component.MetricsManager, components ...component.Core) error {
	t.AuthClient = auth
	t.Permissions = permissions
	t.Settings = settings
	t.Metrics = metrics

	userIds := userids.NewUserIDs(t.Platform.InstallID, t.Platform.MachineID)

	// setup components
	if auth != nil {
		t.Components.Add(auth)
	}

	if permissions != nil {
		t.Components.Add(permissions)
	}

	if settings != nil {
		t.Components.Add(settings)
		settings.AddNotificationTarget(t.Components)
	}

	if metrics != nil {
		t.Components.Add(metrics)
	}

	cohort := cohort.NewTestManager(
		&monetizable.SegmenterMock{
			IsMonetizableReturns: true,
		},
	)
	t.Components.Add(cohort)

	var network component.NetworkManager
	if t.Network != nil {
		network = t.Network
	} else {
		network = NewMockNetworkManager()
	}

	t.Components.Add(network)

	for _, c := range components {
		t.Components.Add(c)
	}

	configuration := config.GetConfiguration(t.Platform)

	t.Components.Initialize(component.InitializerOptions{
		FastnodedURL:      t.Fastnoded.URL,
		Configuration: &configuration,
		AuthClient:    auth,
		License:       auth,
		Cohort:        cohort,
		Permissions:   permissions,
		Settings:      settings,
		Metrics:       metrics,
		Platform:      t.Platform,
		Network:       network,
		UserIDs:       userIds,
	})

	// register HTTP handlers
	t.Components.RegisterHandlers(t.Fastnoded.Router)

	// make sure to empty the (blocking) login / logout channels, if available
	if t.ReadLoginLogoutEvents && auth != nil {
		go t.handleAuthLoop()
	}

	if auth != nil {
		auth.SetTarget(t.Backend.URL)
	}

	if auth != nil {
		// emulates logic in client/http.go
		hasCookie := auth.HasAuthCookie()
		remoteUser, remoteErr := auth.FetchUser(context.Background())
		localUser, localErr := auth.CachedUser()
		switch {
		// If we were able to authenticate remotely, log in. This means the user has a valid
		// session and we were able to fetch the user object remotely
		case remoteErr == nil:
			auth.LoggedInChan() <- remoteUser
			// User will identify in the select loop below via normal login flow

		// If we have an auth cookie and a cached user, treat the user as logged in. This means
		// the user had a valid session before, and a user object was cached during that session.
		// But we currently cannot fetch the user remotely (i.e user is offline)
		case hasCookie && localErr == nil:
			auth.LoggedInChan() <- localUser
			// User will identify in the select loop below via normal login flow
		}
	}

	clienttelemetry.SetCustomTelemetryClient(nil)
	clienttelemetry.SetClientVersion("1.0.0-unit-test")

	return nil
}

// CurrentUser returns the user which is currently logged into fastnoded
func (t *TestClientServer) CurrentUser() (*community.User, error) {
	return t.FastnodedClient.CurrentUser()
}

// CurrentPlan returns the plan for the user which is currently logged into fastnoded
func (t *TestClientServer) CurrentPlan() (*account.PlanResponse, error) {
	return t.FastnodedClient.CurrentPlan()
}

func (t *TestClientServer) handleAuthLoop() {
	for {
		t.HandleAuthEvent()
	}
}

// HandleAuthEvent handles a single pending login or logout events
func (t *TestClientServer) HandleAuthEvent() {
	select {
	case <-t.quitChan:
		return

	case user := <-t.AuthClient.LoggedInChan():
		log.Printf("Login of user %s", user.Name)
		t.AuthClient.SetUser(user)
		t.Components.LoggedIn()
		t.mu.Lock()
		t.fastnodeUser = user
		t.mu.Unlock()

		// make sure that we do not set 0 as userID
		uids := userids.NewUserIDs("", "test-case-machine")
		uids.SetUser(user.ID+1, "", true)
		clienttelemetry.SetUserIDs(uids)

	case <-t.AuthClient.LoggedOutChan():
		log.Printf("User logged out")
		t.Components.LoggedOut()
		t.mu.Lock()
		t.fastnodeUser = nil
		t.mu.Unlock()
	}
}

// GetFilePath returns a sub-path of the Whitelisted base path which is suitable for the current platform
func (t *TestClientServer) GetFilePath(path ...string) string {
	all := append([]string{t.BasePath}, path...)
	return filepath.Join(all...)
}
