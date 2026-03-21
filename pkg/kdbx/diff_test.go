package kdbx

import (
	"testing"
	"time"

	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

func makeEntryAt(title string, t time.Time) Entry {
	entry := makeEntry(title)
	tw := wrappers.TimeWrapper{Time: t}
	entry.Times.LastModificationTime = &tw
	return entry
}

// atTime returns a copy of e with a different timestamp but the same UUID.
func atTime(e Entry, t time.Time) Entry {
	inner := *e.Entry
	tw := wrappers.TimeWrapper{Time: t}
	inner.Times.LastModificationTime = &tw
	return Entry{&inner}
}

func TestDiffDatabases(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	t.Run("empty databases", func(t *testing.T) {
		a := makeDatabase("a.kdbx")
		b := makeDatabase("b.kdbx")

		got := DiffDatabases(a, b)

		if len(got) != 0 {
			t.Errorf("expected empty result, got %v", got)
		}
	})

	t.Run("identical databases", func(t *testing.T) {
		entry := makeEntryAt("GitHub", t0)
		a := makeDatabase("a.kdbx", makeGroup("Passwords", entry))
		b := makeDatabase("b.kdbx", makeGroup("Passwords", entry))

		got := DiffDatabases(a, b)

		if len(got) != 1 {
			t.Fatalf("expected 1 diff, got %d", len(got))
		}
		if got[0].Status != Unchanged {
			t.Errorf("expected Unchanged, got %v", got[0].Status)
		}
	})

	t.Run("entry only in a", func(t *testing.T) {
		entry := makeEntryAt("GitHub", t0)
		a := makeDatabase("a.kdbx", makeGroup("Passwords", entry))
		b := makeDatabase("b.kdbx", makeGroup("Passwords"))

		got := DiffDatabases(a, b)

		if len(got) != 1 {
			t.Fatalf("expected 1 diff, got %d", len(got))
		}
		if got[0].Status != Removed {
			t.Errorf("expected Removed, got %v", got[0].Status)
		}
	})

	t.Run("entry only in b", func(t *testing.T) {
		entry := makeEntryAt("GitHub", t0)
		a := makeDatabase("a.kdbx", makeGroup("Passwords"))
		b := makeDatabase("b.kdbx", makeGroup("Passwords", entry))

		got := DiffDatabases(a, b)

		if len(got) != 1 {
			t.Fatalf("expected 1 diff, got %d", len(got))
		}
		if got[0].Status != Added {
			t.Errorf("expected Added, got %v", got[0].Status)
		}
	})

	t.Run("entry modified", func(t *testing.T) {
		base := makeEntry("GitHub")
		a := makeDatabase("a.kdbx", makeGroup("Passwords", atTime(base, t0)))
		b := makeDatabase("b.kdbx", makeGroup("Passwords", atTime(base, t1)))

		got := DiffDatabases(a, b)

		if len(got) != 1 {
			t.Fatalf("expected 1 diff, got %d", len(got))
		}
		if got[0].Status != Modified {
			t.Errorf("expected Modified, got %v", got[0].Status)
		}
	})

	t.Run("all statuses together", func(t *testing.T) {
		unchanged := makeEntryAt("Unchanged", t0)
		removed := makeEntryAt("Removed", t0)
		added := makeEntryAt("Added", t0)
		base := makeEntry("Modified")

		a := makeDatabase("a.kdbx", makeGroup("G", unchanged, removed, atTime(base, t0)))
		b := makeDatabase("b.kdbx", makeGroup("G", unchanged, added, atTime(base, t1)))

		got := DiffDatabases(a, b)

		if len(got) != 4 {
			t.Fatalf("expected 4 diffs, got %d", len(got))
		}
		statusCount := map[ChangeStatus]int{}
		for _, d := range got {
			statusCount[d.Status]++
		}
		for _, status := range []ChangeStatus{Unchanged, Added, Removed, Modified} {
			if statusCount[status] != 1 {
				t.Errorf("expected exactly 1 of status %v, got %d", status, statusCount[status])
			}
		}
	})

	t.Run("results are sorted by path", func(t *testing.T) {
		entryZ := makeEntryAt("ZEntry", t0)
		entryA := makeEntryAt("AEntry", t0)
		a := makeDatabase("a.kdbx", makeGroup("G", entryZ, entryA))
		b := makeDatabase("b.kdbx", makeGroup("G", entryZ, entryA))

		got := DiffDatabases(a, b)

		for i := 1; i < len(got); i++ {
			if got[i].Path < got[i-1].Path {
				t.Errorf("results not sorted: %v before %v", got[i-1].Path, got[i].Path)
			}
		}
	})

	t.Run("same-path entries are ordered by UUID", func(t *testing.T) {
		// Two entries with the same title (same path) but different UUIDs: one
		// only in A (Removed) and one only in B (Added). The one with the smaller
		// UUID must always come first regardless of map-iteration order.
		low := makeEntryAt("Twin", t0)
		low.UUID = UUID{0x00}
		high := makeEntryAt("Twin", t0)
		high.UUID = UUID{0xff}

		a := makeDatabase("a.kdbx", makeGroup("G", low))
		b := makeDatabase("b.kdbx", makeGroup("G", high))

		got := DiffDatabases(a, b)

		if len(got) != 2 {
			t.Fatalf("expected 2 diffs, got %d", len(got))
		}
		if got[0].UUID != low.UUID {
			t.Errorf("expected low UUID first, got %v", got[0].UUID)
		}
		if got[1].UUID != high.UUID {
			t.Errorf("expected high UUID second, got %v", got[1].UUID)
		}
	})
}
