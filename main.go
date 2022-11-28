package main

import (

	"github.com/shikaan/kpcli/cmd"
)

//go:generate make info

func main() {
	e := cmd.Root.Execute()

	if e != nil {
		println(e.Error())
	}
}
