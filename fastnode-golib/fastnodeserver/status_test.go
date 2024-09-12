package fastnodeserver

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseFastnodeServer(t *testing.T) {
	fastnodeURL, err := ParseFastnodeServerURL("https://tfserving.fastnode.local")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "https://tfserving.fastnode.local:443"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("http://tfserving.fastnode.local")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "http://tfserving.fastnode.local:80"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("http://tfserving.fastnode.local:1234")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "http://tfserving.fastnode.local:1234"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("https://tfserving-1.fastnode.local")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "https://tfserving-1.fastnode.local:443"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("tfserving.fastnode.local:1234")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "http://tfserving.fastnode.local:1234"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("tfserving.fastnode.local")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "https://tfserving.fastnode.local:443"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("tfserving.fastnode.local:1234")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "http://tfserving.fastnode.local:1234"), fastnodeURL)

	fastnodeURL, err = ParseFastnodeServerURL("https://clientname@fastnode.local")
	require.NoError(t, err)
	require.EqualValues(t, mustParse(t, "https://clientname@fastnode.local:443"), fastnodeURL)
}

func Test_Health(t *testing.T) {
	if os.ExpandEnv("CI") != "" {
		t.Skipf("test with network access disabled on CI")
	}

	_, _, err := GetHealth("http://tfserving.khulnasoft.com:8085")
	require.NoError(t, err)

	_, _, err = GetHealth("https://cloud.khulnasoft.com")
	require.NoError(t, err)
}

func mustParse(t *testing.T, value string) *url.URL {
	u, err := url.Parse(value)
	require.NoError(t, err)
	return u
}
