package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/logger"
)

func usage() {
	fmt.Println("Usage: kpcli [COMMAND] [OPTION]... [DATABASE PATH]")
  // TODO: make a prettier help
  // fmt.Println("Open kdbx database located at [DATABASE PATH].")
	fmt.Println("")
	flag.PrintDefaults()
}

func main() {
	l := logger.NewFileLogger(logger.Debug, "kpcli.log")
	defer l.CleanUp()

	flag.Usage = usage
	keyPath := *flag.String("key", "", "Path to the key file to unlock the database")
	password := *flag.String("password", "", "Password to unlock the database")

	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

  if keyPath == "" {
    keyPath = os.Getenv("KPCLI_KEY")
  }

  if password == "" {
    password = os.Getenv("KPCLI_PASSWORD")
  }

	command := flag.Arg(0)
  databasePath := flag.Arg(1)
  var err error

  switch(command) {
  case "list":
    err = List(databasePath, keyPath, password)
  case "copy":
    err = Copy(databasePath, keyPath, password)
  default:
    err = errors.MakeError(fmt.Sprintf("Unrecognized command. Got '%s', expected one of '%s, %s'", command, "list", "open"), "command")
  }
  
  handle(err)
}

func handle(err error) {
  if err != nil {
    println(err.Error())
    os.Exit(1)
  }
}
