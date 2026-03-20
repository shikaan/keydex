package kdbx

import (
	"strings"
	"testing"
)

func TestFormatDiff(t *testing.T) {
	t.Run("headers contain file names", func(t *testing.T) {
		got := FormatDiff("old.kdbx", "new.kdbx", []EntryDiff{})

		if !strings.Contains(got, "--- old.kdbx") {
			t.Errorf("missing --- header, got:\n%s", got)
		}
		if !strings.Contains(got, "+++ new.kdbx") {
			t.Errorf("missing +++ header, got:\n%s", got)
		}
	})

	t.Run("@@ line counts entries in each database", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Unchanged", Status: Unchanged},
			{Path: "/G/Removed", Status: Removed},
			{Path: "/G/Added", Status: Added},
			{Path: "/G/Modified", Status: Modified},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", diffs)

		// a has Unchanged + Removed + Modified = 3
		// b has Unchanged + Added + Modified = 3
		if !strings.Contains(got, "@@ -3 entries +3 entries @@") {
			t.Errorf("unexpected @@ line, got:\n%s", got)
		}
	})

	t.Run("unchanged entries have space prefix", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", []EntryDiff{{Path: "/G/Entry", Status: Unchanged}})

		if !strings.Contains(got, "  /G/Entry") {
			t.Errorf("expected space-prefixed entry, got:\n%s", got)
		}
	})

	t.Run("removed entries have - prefix", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", []EntryDiff{{Path: "/G/Entry", Status: Removed}})

		if !strings.Contains(got, "- /G/Entry") {
			t.Errorf("expected - prefixed entry, got:\n%s", got)
		}
	})

	t.Run("added entries have + prefix", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", []EntryDiff{{Path: "/G/Entry", Status: Added}})

		if !strings.Contains(got, "+ /G/Entry") {
			t.Errorf("expected + prefixed entry, got:\n%s", got)
		}
	})

	t.Run("modified entries have ~ prefix", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", []EntryDiff{{Path: "/G/Entry", Status: Modified}})

		if !strings.Contains(got, "~ /G/Entry") {
			t.Errorf("expected ~ prefixed entry, got:\n%s", got)
		}
	})
}
