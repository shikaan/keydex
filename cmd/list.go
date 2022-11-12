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

	entries := kdbx.GetEntries()

	for _, k := range getSortedKeys(entries) {
		fmt.Println(k)
	}

	return nil
}

func getSortedKeys(entries map[string]*kdbx.Entry) []string {
	keys := []string{}
	for k := range entries {
		keys = append(keys, k)
	}
	less := func(i, j int) bool {
		numberOfSlashesI := len(strings.Split(keys[i], "/"))
		numberOfSlashesJ := len(strings.Split(keys[j], "/"))

    // Sort elements in the same group
		if numberOfSlashesI == numberOfSlashesJ {
			return sort.StringsAreSorted([]string{strings.ToLower(keys[i]), strings.ToLower(keys[j])})
		}

    // Show nested entities close to each other 
		return numberOfSlashesI > numberOfSlashesJ
	}
	sort.Slice(keys, less)
	return keys
}
