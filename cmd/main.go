package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shikaan/kpcli/pages"
	"github.com/shikaan/kpcli/pkg/logger"
)

func usage() {
	fmt.Println("Usage: kpcli [OPTION]... [DATABASE PATH]")
	fmt.Println("Open kdbx database located at [DATABASE PATH].")
	fmt.Println("")
	flag.PrintDefaults()
}

func main() {
	l := logger.NewFileLogger(logger.Debug, "kpcli.log")
	defer l.CleanUp()

	flag.Usage = usage
	keyPath := *flag.String("key", "", "Path to the key file to unlock the database")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	database := flag.Arg(0)

	app.Run(database, keyPath)
}
