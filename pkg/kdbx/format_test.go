package kdbx

import (
	"strings"
	"testing"
)

func TestFormatDiff(t *testing.T) {
	t.Run("returns empty string when there are no changes", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Entry", Status: Unchanged},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", diffs)

		if got != "" {
			t.Errorf("expected empty string, got:\n%s", got)
		}
	})

	t.Run("returns empty string for empty input", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", []EntryDiff{})

		if got != "" {
			t.Errorf("expected empty string, got:\n%s", got)
		}
	})

	t.Run("unchanged entries are not shown", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Unchanged", Status: Unchanged},
			{Path: "/G/Removed", Status: Removed},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", diffs)

		if strings.Contains(got, "/G/Unchanged") {
			t.Errorf("expected unchanged entry to be absent, got:\n%s", got)
		}
	})

	t.Run("headers contain file names", func(t *testing.T) {
		diffs := []EntryDiff{{Path: "/G/Entry", Status: Removed}}

		got := FormatDiff("old.kdbx", "new.kdbx", diffs)

		if !strings.Contains(got, "--- old.kdbx") {
			t.Errorf("missing --- header, got:\n%s", got)
		}
		if !strings.Contains(got, "+++ new.kdbx") {
			t.Errorf("missing +++ header, got:\n%s", got)
		}
	})

	t.Run("@@ line counts all entries in each database", func(t *testing.T) {
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
