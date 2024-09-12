package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/datadeps"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/internal/filesystem"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/internal/indexing"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/fastnodelocal/permissions"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonbatch"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln("unable to get current user:", err)
	}
	indexedDir := usr.HomeDir

	err = datadeps.UseAssetFileMap()
	if err != nil {
		log.Fatalln(err)
	}

	fsManager := filesystem.NewManager(filesystem.Options{
		RootDir:   indexedDir,
		DutyCycle: 0.15,
	})
	fsManager.Initialize(component.InitializerOptions{
		Permissions: permissions.NewManager([]lang.Language{lang.Python}, nil),
	})
	fs := fsManager.FileSystem()

	graph, errc := pythonresource.NewManager(pythonresource.DefaultOptions)
	if err := <-errc; err != nil {
		log.Fatalln(err)
	}

	logf := func(msg string, vals ...interface{}) {
		log.Printf("pythonlocal.select-files: %s", fmt.Sprintf(msg, vals...))
	}

	scanner := bufio.NewScanner(os.Stdin)
	var filename string
	for {
		fmt.Print("Enter filename or 'q' to quit: ")
		scanner.Scan()
		filename = strings.TrimSpace(scanner.Text())
		if filename == "q" {
			break
		}

		logf("starting build with filename: %s", filename)

		ctx := fastnodectx.Background()
		// Select files to index
		ts := time.Now()
		sources, missing, err := pythonbatch.Selector{
			StartPaths:   []string{filename},
			Files:        nil,
			LibraryFiles: fsManager.Files(),
			Opts:         pythonbatch.DefaultOptions.PathSelection,
			Graph:        graph,
			Getter:       indexing.LocalGetter{},
			Logf:         logf,
			FileSystem:   fs,
			Local:        true,
		}.Select(ctx)
		logf("Select took: %s", time.Since(ts))
		if err != nil {
			logf("SelectFiles error: %v", err)
		}
		if len(sources) == 0 {
			if len(missing) > 0 {
				logf("missing file hashes: %d)", len(missing))
			}
			logf("no files selected")
		}
	}
}
