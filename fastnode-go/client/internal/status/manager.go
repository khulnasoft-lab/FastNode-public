package status

import (
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clienttelemetry"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/settings"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/platform"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/localcode"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/codebase"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/response"
	constants "github.com/khulnasoft-lab/fastnode/fastnode-golib/conversion"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/licensing"
)

// NewManager returns a new manager
func NewManager() *Manager {
	return &Manager{
		startTime: time.Now(),
	}
}

// NewTestManager returns a new manager with a dummy uptimeDuration
func NewTestManager() *Manager {
	m := NewManager()
	m.cohort = component.MockCohortManager{
		Convcohort: constants.NoCohort,
	}
	m.settings = settings.NewTestManager()
	// Sensible defaults that can be overridden for specific testing
	m.license = &licensing.MockLicense{
		Plan:    licensing.ProYearly,
		Product: licensing.Pro,
	}
	return m
}

// Manager provides information about the current status of fastnoded
type Manager struct {
	auth        component.AuthClient
	cohort      component.ConversionCohortGetter
	settings    component.SettingsManager
	permissions component.PermissionsManager
	license     licensing.StatusGetter
	platform    *platform.Platform
	startTime   time.Time

	mu              sync.Mutex
	localCodeStatus localcode.StatusResponse
	models          component.IsLoadeder
	nav             component.Validator
}

// Name implements component Core
func (m *Manager) Name() string {
	return "file-status"
}

// Initialize implements component Initializer
func (m *Manager) Initialize(opts component.InitializerOptions) {
	m.auth = opts.AuthClient
	m.permissions = opts.Permissions
	m.platform = opts.Platform
	m.cohort = opts.Cohort
	m.settings = opts.Settings
	m.license = opts.License
}

// EventResponse implements component EventResponser
// it marks the backend as busy after an event was sent by the server to this client
func (m *Manager) EventResponse(root *response.Root) {
	// update local index status from response, if set
	if root != nil && root.LocalIndexStatus != nil {
		m.updateLocalCodeStatus(*root.LocalIndexStatus)
	}
}

// RegisterHandlers implements component Handler
func (m *Manager) RegisterHandlers(mux *mux.Router) {
	mux.HandleFunc("/clientapi/status", m.handleStatus).Methods("GET")
}

// SetNav implements StatusManager
func (m *Manager) SetNav(nav component.Validator) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nav = nav
}

func (m *Manager) navValidate(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.nav == nil {
		// SetNav has not been called means Fastnode is initializing
		return codebase.ErrProjectStillIndexing
	}
	return m.nav.Validate(path)
}

// SetModels implements StatusManager
func (m *Manager) SetModels(lm component.IsLoadeder) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.models = lm
}

func (m *Manager) isLoaded(fext string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.models == nil {
		return false
	}

	return m.models.IsLoaded(fext)
}

func (m *Manager) updateLocalCodeStatus(status localcode.StatusResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()

	hadLoaded := len(m.localCodeStatus.Indices) > 0
	loadedNow := len(status.Indices) > 0

	switch {
	case hadLoaded && loadedNow:
		// Already had an index
	case !hadLoaded && loadedNow:
		// Index loaded when previously one was not
		clienttelemetry.FastnodeTelemetry("Local Index Added", map[string]interface{}{
			"local_code_status": status,
		})
	case hadLoaded && !loadedNow:
		// Index went away
		clienttelemetry.FastnodeTelemetry("Local Index Removed", map[string]interface{}{
			"local_code_status": status,
		})
	}

	m.localCodeStatus = status
}
