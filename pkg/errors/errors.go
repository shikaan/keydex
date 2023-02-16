package errors

import (
	"fmt"

	"github.com/shikaan/keydex/pkg/info"
)

func MakeError(msg string, namespace string) error {
	return fmt.Errorf("%s(%s): %s", info.NAME, namespace, msg)
}
