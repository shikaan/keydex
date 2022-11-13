package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shikaan/kpcli/pkg/kdbx"
)

func List(database, key, passphrase string) error {
	kdbx, err := kdbx.NewUnlocked(database, passphrase)
	if err != nil {
		return err
	}

	entries := kdbx.GetEntryPaths()

	for _, k := range getSortedKeys(entries) {
		fmt.Println(k)
	}

	return nil
}

func getSortedKeys(entries []kdbx.EntryPath) []kdbx.EntryPath {
	less := func(i, j int) bool {
		numberOfSlashesI := len(strings.Split(entries[i], "/"))
		numberOfSlashesJ := len(strings.Split(entries[j], "/"))

    // Sort elements in the same group
		if numberOfSlashesI == numberOfSlashesJ {
			return sort.StringsAreSorted([]string{strings.ToLower(entries[i]), strings.ToLower(entries[j])})
		}

    // Show nested entities close to each other 
		return numberOfSlashesI > numberOfSlashesJ
	}
	sort.Slice(entries, less)
	return entries
}
