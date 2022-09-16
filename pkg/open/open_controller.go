package open

import (
	"fmt"
	"syscall"

	"golang.org/x/term"

	"github.com/shikaan/kpcli/pkg/kdbx"
)

func Open(database, keyPath string) {
	fmt.Println("Insert password: ")

	pwd, err := term.ReadPassword(int(syscall.Stdin))
	handleError(err)

	kdbx, err := kdbx.New(database)
	handleError(err)

	err = kdbx.Unlock(string(pwd))
	handleError(err)

	Render(*kdbx)
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}
