package localcode

import "fmt"

// ArtifactTempDir returns a root directory to place temporary artifacts associated with artifact UUID provided.
func ArtifactTempDir(id string) string {
	return fmt.Sprintf("/var/fastnode/tmp/artifacts/%s", id)
}
