package clipboard

import (
	"github.com/atotto/clipboard"
	"github.com/shikaan/kpcli/pkg/errors"
)

func Write(msg string) error {
	err := clipboard.WriteAll(msg)

	if err != nil {
		return errors.MakeError(err.Error(), "clipboard")
	}

	return nil
}
