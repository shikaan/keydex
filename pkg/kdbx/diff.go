package kdbx

import (
	"slices"
	"strings"
	"time"
)

type ChangeStatus int

const (
	Unchanged ChangeStatus = iota
	Added
	Removed
	Modified
)

type EntryDiff struct {
	UUID   UUID
	Path   EntityPath
	Status ChangeStatus
}

type entryRecord struct {
	path EntityPath
	time time.Time
}

func DiffDatabases(a, b *Database) []EntryDiff {
	aMap := makeRecordMap(a)
	bMap := makeRecordMap(b)

	seen := map[UUID]struct{}{}
	for uuid := range aMap {
		seen[uuid] = struct{}{}
	}
	for uuid := range bMap {
		seen[uuid] = struct{}{}
	}

	result := []EntryDiff{}
	for uuid := range seen {
		aRec, inA := aMap[uuid]
		bRec, inB := bMap[uuid]

		var status ChangeStatus
		var path EntityPath
		switch {
		case inA && !inB:
			status = Removed
			path = aRec.path
		case !inA && inB:
			status = Added
			path = bRec.path
		case aRec.time.Equal(bRec.time):
			status = Unchanged
			path = bRec.path
		default:
			status = Modified
			path = bRec.path
		}
		result = append(result, EntryDiff{UUID: uuid, Path: path, Status: status})
	}

	slices.SortFunc(result, func(a, b EntryDiff) int {
		return strings.Compare(a.Path, b.Path)
	})

	return result
}

func makeRecordMap(db *Database) map[UUID]entryRecord {
	result := map[UUID]entryRecord{}
	for _, g := range db.Content.Root.Groups {
		collectRecords(g, PATH_SEPARATOR, result)
	}
	return result
}

func collectRecords(g Group, prefix string, out map[UUID]entryRecord) {
	groupPrefix := makeGroupPrefix(prefix, g)

	for _, sub := range g.Groups {
		collectRecords(sub, groupPrefix, out)
	}

	for _, entry := range g.Entries {
		var t time.Time
		if entry.Times.LastModificationTime != nil {
			t = entry.Times.LastModificationTime.Time
		}
		out[entry.UUID] = entryRecord{
			path: makeEntryPath(groupPrefix, entry),
			time: t,
		}
	}
}
