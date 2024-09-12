package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/navigation/localpath"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	khulnasoft-lab := localpath.Absolute(filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "khulnasoft-lab", "khulnasoft-lab"))
	s, err := NewStorage(StorageOptions{})
	require.NoError(t, err)

	var noCache []Commit
	for i := 0; i < 3; i++ {
		repo, err := Open(khulnasoft-lab, DefaultComputedCommitsLimit, s)
		require.NoError(t, err)

		var batch []Commit
		for j := 0; j < 10; j++ {
			commit, err := repo.Next(fastnodectx.Background())
			require.NoError(t, err)
			batch = append(batch, commit)
			if i == 0 {
				noCache = append(noCache, commit)
				continue
			}
			require.Equal(t, noCache[j], commit)
		}

		err = repo.Save(s)
		require.NoError(t, err)
	}
}
