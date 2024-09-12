// +build standalone,!linux

package statusicon

import "github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/updates"

// MockManager implements component.Core
type MockManager struct{}

// Name implements interface Core
func (MockManager) Name() string {
	return "statusicon_mock"
}

// NewManager returns a new statusicon component
func NewManager(_ updates.Manager) MockManager {
	return MockManager{}
}
