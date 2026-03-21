package kdbx

import (
	"fmt"
	"strings"
	"time"

	"github.com/tobischo/gokeepasslib/v3"
)

func formatGroupPrefix(prefix string, g Group) string {
	return prefix + g.Name + PATH_SEPARATOR
}

func formatEntryPath(groupPrefix string, entry gokeepasslib.Entry) EntityPath {
	title := entry.GetTitle()
	if title == "" {
		title = "(UNKNOWN)"
	}
	return groupPrefix + sanitizePathPortion(title)
}

const timestampLayout = "2006-01-02 15:04:05"

func FormatDiff(nameA, nameB string, timeA, timeB time.Time, diffs []EntryDiff) string {
	var countA, countB int
	var changed []EntryDiff

	for _, d := range diffs {
		switch d.Status {
		case Unchanged:
			countA++
			countB++
		case Removed:
			countA++
			changed = append(changed, d)
		case Added:
			countB++
			changed = append(changed, d)
		case Modified:
			countA++
			countB++
			changed = append(changed, d)
		}
	}

	if len(changed) == 0 {
		return ""
	}

	startA, startB := 1, 1
	if countA == 0 {
		startA = 0
	}
	if countB == 0 {
		startB = 0
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\t%s\n", nameA, timeA.Format(timestampLayout))
	fmt.Fprintf(&sb, "+++ %s\t%s\n", nameB, timeB.Format(timestampLayout))
	fmt.Fprintf(&sb, "@@ -%d,%d +%d,%d @@\n", startA, countA, startB, countB)

	for _, d := range changed {
		switch d.Status {
		case Removed:
			fmt.Fprintf(&sb, "-%s\n", d.Path)
		case Added:
			fmt.Fprintf(&sb, "+%s\n", d.Path)
		case Modified:
			fmt.Fprintf(&sb, "-%s\n", d.Path)
			fmt.Fprintf(&sb, "+%s\n", d.Path)
		}
	}

	return sb.String()
}
