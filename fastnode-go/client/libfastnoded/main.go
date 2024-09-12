package main

/*
#cgo darwin CFLAGS: -mmacosx-version-min=10.11
#cgo darwin LDFLAGS: -mmacosx-version-min=10.11
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"sync"

	clientpkg "github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/client"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/clientapp"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/updates/liveupdates"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/sidebar"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/visibility"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/rollbar"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/throttle"
)

var (
	mu     sync.Mutex
	client *clientpkg.Client
)

//export fastnodeInitialize
func fastnodeInitialize() (ret bool) {
	defer panicRecoveryBool(&ret)

	if err := throttle.SetLowPriority(); err != nil {
		log.Printf("failed to set low process priority: %s", err)
	}

	mu.Lock()
	defer mu.Unlock()

	var err error
	target, err := liveupdates.UpdateTarget()
	if err != nil {
		clientapp.Alert(err)
		rollbar.Critical(fmt.Errorf("fastnodeInitialize: error initializing libfastnoded: %v", err))
		log.Println("fastnodeInitialize: error initializing libfastnoded:", err)
		return false
	}

	client, err = clientapp.Start(&clientpkg.Options{Updater: liveupdates.NewManager(target)})
	if err != nil {
		clientapp.Alert(err)
		rollbar.Critical(fmt.Errorf("fastnodeInitialize: error initializing libfastnoded: %v", err))
		log.Println("fastnodeInitialize: error initializing libfastnoded:", err)
		return false
	}

	return true
}

//export fastnodeConnect
func fastnodeConnect() (ret bool) {
	defer panicRecoveryBool(&ret)

	mu.Lock()
	defer mu.Unlock()
	if client == nil {
		log.Println("fastnodeConnectToHost: run FastnodeInitialize before connecting")
		return false
	}

	go func() {
		err := client.Connect(client.Settings.Server())
		if err != nil {
			log.Printf("libfastnoded failed to connect: %v", err)
			// TODO: have the outer function wait for confirmation that fastnode is connected
			// before returning, and then on failure have it return false.
		}
	}()

	return true
}

//export fastnodeSetEnv
func fastnodeSetEnv(key, value *C.char) {
	defer panicRecovery()

	keyStr := C.GoString(key)
	valueStr := C.GoString(value)
	os.Setenv(keyStr, valueStr)
}

//export fastnodeTrackSidebarVisibility
// fastnodeTrackSidebarVisibility stores state about whether the sidebar was visible at the
// time the method is invoked
func fastnodeTrackSidebarVisibility() {
	defer panicRecovery()

	mu.Lock()
	defer mu.Unlock()
	if client == nil {
		log.Println("fastnodeTrackSidebar: run FastnodeInitialize before connecting")
		return
	}

	sidebar.SetRestartIfPreviouslyVisible(visibility.RecentlyVisible())
}

//export fastnodeStopSidebar
// fastnodeStopSidebar stops the sidebar application if it is running
func fastnodeStopSidebar() {
	defer panicRecovery()
	sidebar.Stop()
}

//export fastnodeCheckForUpdates
// fastnodeCheckForUpdates is called when the user explicitly checks for updates
func fastnodeCheckForUpdates(showModal bool) {
	defer panicRecovery()

	mu.Lock()
	defer mu.Unlock()
	if client == nil {
		log.Println("fastnodeCheckForUpdates: run FastnodeInitialize before connecting")
		return
	}

	client.Updater.CheckForUpdates(showModal)
}

//export fastnodeUpdateReady
// fastnodeUpdateReady returns true if an update has been downloaded and is waiting to
// be installed when the process terminates.
func fastnodeUpdateReady() bool {
	defer panicRecovery()

	mu.Lock()
	defer mu.Unlock()
	if client == nil {
		log.Println("fastnodeUpdateReady: run FastnodeInitialize before connecting")
		return false
	}

	return client.Updater.UpdateReady()
}

func panicRecoveryBool(ret *bool) {
	if err := recover(); err != nil {
		rollbar.PanicRecovery(err)
		*ret = false
	}
}

func panicRecovery() {
	if err := recover(); err != nil {
		rollbar.PanicRecovery(err)
	}
}

func main() {
	// Required by CGO. Does nothing. See exported API above
}
