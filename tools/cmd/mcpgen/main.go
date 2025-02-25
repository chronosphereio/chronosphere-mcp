package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chronosphereio/chronoctl-core/tools/pkg/clispec"
	"mvdan.cc/gofumpt/format"
)

var (
	specPath        = flag.String("spec", "", "location of spec file")
	targetDir       = flag.String("target", "./cli", "location of target directory")
	packageName     = flag.String("pkg", "", "name of the generated package")
	allowedEntities = flag.String("allowed-entities", "", "comma seperated list of explicitly allowed entities. If omitted, all entities are generated")
)

func main() {
	flag.Parse()

	if strings.Contains(*specPath, "unstable") && *allowedEntities == "" {
		log.Fatal("must allow list generated unstable entities using --allowed-entities")
	}

	spec, err := clispec.Parse(*specPath)
	if err != nil {
		log.Fatal("could not create spec handler: ", err)
	}

	files, err := generate(*packageName, spec)
	if err != nil {
		log.Fatal("could not generate: ", err)
	}
	for _, file := range files {
		if err := writeGoFile(filepath.Join(*targetDir, file.name+".gen.go"), []byte(file.content)); err != nil {
			log.Fatal("could not write: ", err)
		}
	}
}

func writeGoFile(name string, contents []byte) (retErr error) {
	writeContents, err := format.Source(contents, format.Options{})
	if err != nil {
		// Write the unformatted content for easier debugging.
		writeContents = contents
		defer func() {
			if retErr == nil {
				retErr = fmt.Errorf("format source file %q failed (unformatted content written): %v", name, err)
			}
		}()
	}
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		log.Fatalf("could not create dir %v: %v", filepath.Dir(name), err)
	}
	return os.WriteFile(name, writeContents, 0o666)
}
