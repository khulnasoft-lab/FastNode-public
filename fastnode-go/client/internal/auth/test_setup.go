package auth

import (
	"fmt"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/mockserver"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/visibility"
)

// SetupWithAuthDefaults configures a default auth client and permissions manager
// both and all other components passed as argument are registered in the Fastnoded mock server as request Handlers
func SetupWithAuthDefaults(t *mockserver.TestClientServer, components ...component.Core) error {
	visibility.Clear()

	authClient := NewTestClient(300 * time.Millisecond)

	err := t.SetupWithCustomAuthClient(authClient, components...)
	if err != nil {
		return err
	}

	// make sure that all HTTP responses were closed by our components
	t.AddCleanupAction(func() {
		unclosed := authClient.getOpenConnections()
		if unclosed > 0 {
			panic(fmt.Sprintf("auth client has unclosed connections: %d", unclosed))
		}
	})

	authClient.SetTarget(t.Backend.URL)
	return nil
}
