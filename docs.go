//go:build exclude

package main

import (
	"archive/tar"
	"log"
	"os"
	"path/filepath"

	"github.com/shikaan/keydex/cmd"
	"github.com/shikaan/keydex/pkg/info"
	"github.com/spf13/cobra/doc"
)

const manpath = "./.build"
const docspath = "./docs"

func main() {
	if err := os.MkdirAll(docspath, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(manpath, 0o755); err != nil {
		log.Fatal(err)
	}

	if err := doc.GenMarkdownTree(cmd.Root, docspath); err != nil {
		log.Fatal(err)
	}

	hdr := &doc.GenManHeader{Title: info.NAME, Section: "1"}
	if err := doc.GenManTree(cmd.Root, hdr, manpath); err != nil {
		log.Fatal(err)
	}

	archive, err := os.Create(filepath.Join(manpath, info.NAME+".1.tar"))
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	tw := tar.NewWriter(archive)
	defer tw.Close()

	pages, err := filepath.Glob(filepath.Join(manpath, "*.1"))
	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		data, err := os.ReadFile(page)
		if err != nil {
			log.Fatal(err)
		}

		hdr := &tar.Header{
			Name: filepath.Base(page),
			Mode: 0o644,
			Size: int64(len(data)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}

		if _, err := tw.Write(data); err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(page); err != nil {
			log.Fatal(err)
		}
	}
}
