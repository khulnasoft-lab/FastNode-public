package test

import (
	"runtime"
	"testing"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clienttelemetry"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/metrics/livemetrics"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/response"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EditorEvents(t *testing.T) {
	testEnv, err := clientapp.StartEmptyTestEnvironment()
	require.NoError(t, err)
	defer testEnv.Close()

	m := testEnv.Fastnoded.Metrics.(*livemetrics.Manager)

	tracker := telemetry.MockClient{}
	clienttelemetry.SetCustomTelemetryClient(&tracker)

	// login to enable tracking
	_, err = testEnv.FastnodedClient.SendLoginRequest("user@example.com", "secret", true)
	require.NoError(t, err)

	flat := m.Completions().ReadAndFlatten(false, nil)
	assert.EqualValues(t, nil, flat["completions_num_selected"])
	assert.EqualValues(t, nil, flat["completions_requested"])

	// emulate edit "conten" -> "content" by completion
	_, err = testEnv.FastnodedClient.PostEditEvent("test_client", "file.py", "conten", 6)
	require.NoError(t, err)

	m.Completions().Get(lang.Python).Requested()
	flat = m.Completions().ReadAndFlatten(false, nil)
	assert.EqualValues(t, 1, flat["completions_requested"])

	m.Completions().Get(lang.Python).ReturnedCompat("content", 6, []response.EditorCompletion{{Insert: "content"}}, time.Now().Add(-15*time.Millisecond))
	flat = m.Completions().ReadAndFlatten(false, nil)
	assert.EqualValues(t, 1, flat["completions_shown"])

	_, err = testEnv.FastnodedClient.PostEditEvent("test_client", "file.py", "content", 7)
	require.NoError(t, err)

	if runtime.GOOS == "darwin" {
		// 0 tracked messages for library walk expected on mac
		require.EqualValues(t, 0, len(tracker.TrackedFilteredByEvent("Background Library Walk Completed")), "Messages: %s", tracker.Tracked())
	}
}
