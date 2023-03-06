//go:build exclude

package main

import (
	"flag"
	"os"
	"os/exec"
	"strings"
)

const INFO_TEMPLATE = "pkg/info/info.tmpl"
const INFO_DESTINATION = "pkg/info/info.go"

func revision() string {
	output, e := exec.Command("git", "log", "-n1", "--pretty=%h").Output()
	if e != nil {
		panic(e)
	}
	return strings.TrimSpace(string(output[:]))
}

func version() string {
	output, _ := exec.Command("git", "describe", "--abbrev=0", "--tags").Output()
	result := strings.TrimSpace(string(output[:]))

	if result == "" {
		return "dev"
	}

	return result
}

func main() {
	revision := flag.String("revision", revision(), "Revision to identify the build, usually a git sha")
	version := flag.String("version", version(), "Version to identified last published version of the package, usually a git tag")
	name := flag.String("name", "keydex", "Name of this executable")

	flag.Parse()

	buffer, err := os.ReadFile(INFO_TEMPLATE)

	if err != nil {
		panic(err)
	}

  // This is needed for local development where env.VERSION might be nil
  versionWithDefault := *version
  if versionWithDefault == "" {
    versionWithDefault = "dev"
  } 

	result := strings.ReplaceAll(string(buffer[:]), "_REVISION_", *revision)
	result = strings.ReplaceAll(result, "_VERSION_", versionWithDefault)
	result = strings.ReplaceAll(result, "_NAME_", *name)

	if err := os.WriteFile(INFO_DESTINATION, []byte(result), 0666); err != nil {
		panic(err)
	}
}
