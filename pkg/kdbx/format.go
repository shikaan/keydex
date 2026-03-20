package kdbx

import (
	"fmt"
	"strings"

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

func FormatDiff(nameA, nameB string, diffs []EntryDiff) string {
	var countA, countB int
	for _, d := range diffs {
		switch d.Status {
		case Unchanged:
			countA++
			countB++
		case Removed:
			countA++
		case Added:
			countB++
		case Modified:
			countA++
			countB++
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\n", nameA)
	fmt.Fprintf(&sb, "+++ %s\n", nameB)
	fmt.Fprintf(&sb, "@@ -%d entries +%d entries @@\n", countA, countB)

	for _, d := range diffs {
		var prefix string
		switch d.Status {
		case Unchanged:
			prefix = " "
		case Removed:
			prefix = "-"
		case Added:
			prefix = "+"
		case Modified:
			prefix = "~"
		}
		fmt.Fprintf(&sb, "%s %s\n", prefix, d.Path)
	}

	return sb.String()
}
