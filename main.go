package main

import (

	"github.com/shikaan/kpcli/cmd"
)

//go:generate make docs

func main() {
	e := cmd.Root.Execute()

	if e != nil {
		println(e.Error())
	}
}
