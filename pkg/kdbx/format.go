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

	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\n", nameA)
	fmt.Fprintf(&sb, "+++ %s\n", nameB)
	fmt.Fprintf(&sb, "@@ -%d entries +%d entries @@\n", countA, countB)

	for _, d := range changed {
		var prefix string
		switch d.Status {
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
