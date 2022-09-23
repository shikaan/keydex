package errors

import (
	"errors"
	"fmt"
)

const AppName = "kpcli"

func MakeError(msg string, namespace string) error {
	return errors.New(fmt.Sprintf("%s(%s): %s", AppName, namespace, msg))
}
