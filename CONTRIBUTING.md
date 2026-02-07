Contributing
---

Practical guidelines to contribute to this project.

## Building from source

In order to build the app from source, you can just type:

```sh
# Get all the dependencies
go get ./...

# Generate app metadata (e.g., revision, version from git history)
make info

# Generate the binary for the current architecture and os
make build
```

## Run locally

The easiest way is to use something like `direnv` and define a `.envrc` file like

```sh
export KEYDEX_PASSPHRASE="my-password"
export KEYDEX_DATABASE="test.kdbx"
```

The `.envrc` and the `test.kdbx` files are gitignored already.

## Sharing a testing build

If you need to share a testing build, maybe because you are working on a complicated
bug, you can simply add your branch name to the [configuration](./.github/workflows/build.yml)
file.

It will create a new Github Release tagged with the branch name.

## Manual testing

Most flows are covered by automated tests in `test/tui_test.go` and
`test/cli_test.go`. The smoke tests below focus on side effects that
cannot be verified automatically (clipboard, file persistence,
interactive input).

<details>
<summary>Clipboard</summary>

**Copy entry password**
1. Open Entry List (^P), select an entry
2. Reveal password (^R), then copy (^C)
   - **Expected:** Password is in system clipboard

**Copy password field**
1. `keydex copy <archive> <ref>`
   - **Expected:** Password is copied to system clipboard

**Copy non-password field**
1. `keydex copy --field UserName <archive> <ref>`
   - **Expected:** Field value is copied to system clipboard

</details>

<details>
<summary>File persistence</summary>

**Save persists to disk**
1. Open an entry, modify a field, save (^O → Y)
2. Close keydex, reopen the same archive
   - **Expected:** Changes are still present

**Delete persists to disk**
1. Open an entry, delete (^D → Y)
2. Close keydex, reopen the same archive
   - **Expected:** Entry no longer exists

</details>

<details>
<summary>Interactive password prompt</summary>

1. Run keydex without `KEYDEX_PASSPHRASE` set
   - **Expected:** Prompted for password
2. Type the correct password
   - **Expected:** Database opens successfully

</details>

<details>
<summary>Piping</summary>

1. `keydex list <archive> | grep GitHub`
   - **Expected:** Only matching entries printed
2. `echo "/TestDB/Coding/GitHub" | keydex copy <archive>`
   - **Expected:** Password copied to clipboard from piped reference
3. `keydex list | fzf | keydex copy`
   - **Expected:** Browse list with fzf and copy to clipboard

</details>

<details>
<summary>Smoke tests</summary>

Quick walkthrough of the most important flows to sanity-check a
release. These are all covered by automated tests too, but a human
pair of eyes helps catch visual glitches and UX regressions.

**Create, save, and verify a new entry**
1. Open keydex, create entry (^N)
2. Fill in fields, save (^O → Y)
3. Open Entry List (^P)
   - **Expected:** New entry appears, no [MODIFIED] banner

**Edit an existing entry**
1. Open Entry List (^P), select an entry
2. Change a field, save (^O → Y)
   - **Expected:** Notification shown, [MODIFIED] disappears, UI reflects changes

**Delete an entry**
1. Open Entry List (^P), select an entry
2. Delete (^D → Y)
   - **Expected:** Entry removed from list

**Read-only mode**
1. Open archive with `--read-only`
2. Attempt to edit, save, or delete
   - **Expected:** Each action shows a read-only notification

**List entries (CLI)**
1. `keydex list <archive>`
   - **Expected:** All entries printed to stdout

**Error handling**
1. `keydex open non-existent.kdbx`
   - **Expected:** `"no such file or directory"` on stderr
2. Open archive with wrong password
   - **Expected:** `"Wrong password?"` on stderr

</details>
