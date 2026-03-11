package clipboard

import (
	"github.com/atotto/clipboard"
	"github.com/shikaan/keydex/pkg/errors"
)

func Write(msg string) error {
	err := clipboard.WriteAll(msg)

	if err != nil {
		return errors.MakeError("Clipboard error: "+err.Error(), "clipboard")
	}

	return nil
}
