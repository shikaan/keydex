package logger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		s := &strings.Builder{}
		l := New(Error, s)
		l.Error("error")
		l.Info("info")
		l.Debug("debug")

		assert.Contains(t, s.String(), "error")
		assert.NotContains(t, s.String(), "info")
		assert.NotContains(t, s.String(), "debug")
	})

	t.Run("info", func(t *testing.T) {
		s := &strings.Builder{}
		l := New(Info, s)
		l.Error("error")
		l.Info("info")
		l.Debug("debug")

		assert.Contains(t, s.String(), "error")
		assert.Contains(t, s.String(), "info")
		assert.NotContains(t, s.String(), "debug")
	})

	t.Run("debug", func(t *testing.T) {
		s := &strings.Builder{}
		l := New(Debug, s)
		l.Error("error")
		l.Info("info")
		l.Debug("debug")

		assert.Contains(t, s.String(), "error")
		assert.Contains(t, s.String(), "info")
		assert.Contains(t, s.String(), "debug")
	})
}
