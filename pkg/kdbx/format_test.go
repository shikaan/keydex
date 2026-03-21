package kdbx

import (
	"strings"
	"testing"
	"time"
)

func TestFormatDiff(t *testing.T) {
	t.Run("returns empty string when there are no changes", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Entry", Status: Unchanged},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, diffs)

		if got != "" {
			t.Errorf("expected empty string, got:\n%s", got)
		}
	})

	t.Run("returns empty string for empty input", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, []EntryDiff{})

		if got != "" {
			t.Errorf("expected empty string, got:\n%s", got)
		}
	})

	t.Run("unchanged entries are not shown", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Unchanged", Status: Unchanged},
			{Path: "/G/Removed", Status: Removed},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, diffs)

		if strings.Contains(got, "/G/Unchanged") {
			t.Errorf("expected unchanged entry to be absent, got:\n%s", got)
		}
	})

	t.Run("headers contain file names", func(t *testing.T) {
		diffs := []EntryDiff{{Path: "/G/Entry", Status: Removed}}

		got := FormatDiff("old.kdbx", "new.kdbx", time.Time{}, time.Time{}, diffs)

		if !strings.Contains(got, "--- old.kdbx") {
			t.Errorf("missing --- header, got:\n%s", got)
		}
		if !strings.Contains(got, "+++ new.kdbx") {
			t.Errorf("missing +++ header, got:\n%s", got)
		}
	})

	t.Run("headers contain timestamps", func(t *testing.T) {
		tA := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		tB := time.Date(2024, 6, 20, 18, 45, 0, 123456789, time.UTC)
		diffs := []EntryDiff{{Path: "/G/Entry", Status: Removed}}

		got := FormatDiff("a.kdbx", "b.kdbx", tA, tB, diffs)

		wantA := "--- a.kdbx\t2024-01-15 10:30:00\n"
		wantB := "+++ b.kdbx\t2024-06-20 18:45:00\n"
		if !strings.Contains(got, wantA) {
			t.Errorf("missing or malformed --- header, got:\n%s", got)
		}
		if !strings.Contains(got, wantB) {
			t.Errorf("missing or malformed +++ header, got:\n%s", got)
		}
	})

	t.Run("@@ line follows POSIX unified diff format", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Unchanged", Status: Unchanged},
			{Path: "/G/Removed", Status: Removed},
			{Path: "/G/Added", Status: Added},
			{Path: "/G/Modified", Status: Modified},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, diffs)

		// a has Removed + Modified = 2 (Unchanged is not emitted)
		// b has Added + Modified = 2 (Unchanged is not emitted)
		if !strings.Contains(got, "@@ -1,2 +1,2 @@") {
			t.Errorf("unexpected @@ line, got:\n%s", got)
		}
	})

	t.Run("@@ start line is 0 when a is empty (all entries added)", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Entry", Status: Added},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, diffs)

		if !strings.Contains(got, "@@ -0,0 +1,1 @@") {
			t.Errorf("expected @@ -0,0 +1,1 @@, got:\n%s", got)
		}
	})

	t.Run("@@ start line is 0 when b is empty (all entries removed)", func(t *testing.T) {
		diffs := []EntryDiff{
			{Path: "/G/Entry", Status: Removed},
		}

		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, diffs)

		if !strings.Contains(got, "@@ -1,1 +0,0 @@") {
			t.Errorf("expected @@ -1,1 +0,0 @@, got:\n%s", got)
		}
	})

	t.Run("removed entries have - prefix", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, []EntryDiff{{Path: "/G/Entry", Status: Removed}})

		if !strings.Contains(got, "-/G/Entry") {
			t.Errorf("expected - prefixed entry, got:\n%s", got)
		}
	})

	t.Run("added entries have + prefix", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, []EntryDiff{{Path: "/G/Entry", Status: Added}})

		if !strings.Contains(got, "+/G/Entry") {
			t.Errorf("expected + prefixed entry, got:\n%s", got)
		}
	})

	t.Run("modified entries are shown as removal followed by addition", func(t *testing.T) {
		got := FormatDiff("a.kdbx", "b.kdbx", time.Time{}, time.Time{}, []EntryDiff{{Path: "/G/Entry", Status: Modified}})

		lines := strings.Split(strings.TrimSpace(got), "\n")
		var entryLines []string
		for _, l := range lines {
			if strings.HasPrefix(l, "-") || strings.HasPrefix(l, "+") {
				if !strings.HasPrefix(l, "---") && !strings.HasPrefix(l, "+++") {
					entryLines = append(entryLines, l)
				}
			}
		}

		if len(entryLines) != 2 {
			t.Fatalf("expected 2 entry lines for modified, got %d:\n%s", len(entryLines), got)
		}
		if entryLines[0] != "-/G/Entry" {
			t.Errorf("expected first line to be removal, got: %s", entryLines[0])
		}
		if entryLines[1] != "+/G/Entry" {
			t.Errorf("expected second line to be addition, got: %s", entryLines[1])
		}
	})
}
