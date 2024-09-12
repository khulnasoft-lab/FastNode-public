package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/refactor/importgraph"
)

const (
	khulnasoft-lab = "github.com/khulnasoft-lab/fastnode"
)

func main() {
	var target string
	flag.StringVar(&target, "target", "", "target to check against")
	flag.Parse()

	// Read in changed files
	scanner := bufio.NewScanner(os.Stdin)
	var changedFiles []string
	for scanner.Scan() {
		changedFiles = append(changedFiles, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// Build the import graph
	forward, _, _ := importgraph.Build(&build.Default)
	reachable := forward.Search(target)

	// Filter out only khulnasoft-lab imports
	khulnasoft-labImports := make(map[string]bool)
	for pkg, imported := range reachable {
		if imported && strings.HasPrefix(pkg, khulnasoft-lab) {
			khulnasoft-labImports[pkg] = imported
		}
	}

	// Check to see if any of the changed files are in khulnasoft-lab imports for the target
	// NOTE: Assumes changedFiles paths are relative to khulnasoft-lab.
	var changed bool
	for _, cf := range changedFiles {
		filename := path.Join(khulnasoft-lab, cf)
		pkg := strings.TrimSuffix(path.Dir(filename), "/")
		if _, exists := khulnasoft-labImports[pkg]; exists {
			fmt.Println(pkg, "->", filename)
			changed = true
		}
	}

	if changed {
		os.Exit(1)
	}
}
