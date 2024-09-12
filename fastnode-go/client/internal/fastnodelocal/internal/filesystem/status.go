package filesystem

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/status"

var (
	section = status.NewSection("client/internal/fastnodelocal/internal/filesystem")

	syncDirCount = section.CounterDistribution("Sync directories (when enabled)")

	// files
	filesCount  = section.Counter("Calls to fs.Files")
	storeCount  = section.Counter("Calls to files.Store")
	deleteCount = section.Counter("Calls to files.Delete")

	// walker
	// TODO(hrysoula): add more walker stats
	walkStatCount = section.CounterDistribution("Calls to os.Lstat during walk")

	// watcher
	eventsPerGroup = section.CounterDistribution("Events per group")
)
