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

Before rolling out a big change, it's good measure to test out the 
most important flows of keydex.

For reference, most of them are documented here.

<details>
<summary>open: no read only, no selected entry</summary>

**Test 1: Create and save new entry**
1. Create Entry in Welcome Screen (^N)
   - **Expected:** Password has no ending `==`, is random, and is at least 16 characters
2. Save (^O)
3. Confirm (Y)
4. Open Entry List (^P)
   - **Expected:** Entry appears in list of entries

**Test 2: Create and dismiss new entry**
1. Create Entry in Welcome Screen (^N)
2. Save (^O)
3. Dismiss (N)
4. Open Entry List (^P)
   - **Expected:** Entry does NOT appear in list of entries

**Test 3: Create and cancel new entry**
1. Create Entry in Welcome Screen (^N)
2. Cancel (ESC)
   - **Expected:** Returns to entry list
3. Open Entry List (^P)
   - **Expected:** Entry does NOT appear in list

**Test 4: View and copy entry password**
1. Open Entry List (^P)
2. Select an Entry
3. Reveal password (^R)
   - **Expected:** Password is readable
4. Copy password (^C)
   - **Expected:** Password is in clipboard

**Test 5: Edit and cancel non-hidden and hidden fields**
1. Open Entry List (^P)
2. Select an Entry
3. Update non-hidden field
   - **Expected:** Field changes
4. Cancel (ESC)
   - **Expected:** Entry returns to previous state
5. Update hidden field
   - **Expected:** Error displayed
6. Reveal field (^R)
   - **Expected:** Password can be updated
7. Update password
8. Cancel (ESC)
   - **Expected:** Entry returns to previous state
9. Save (^O)
   - **Expected:** Nothing happens (no changes to save)

**Test 6: Update entry title**
1. Open Entry List (^P)
2. Select an Entry
3. Update title field
   - **Expected:** Field changes
4. Save (^O)
   - **Expected:** Notification displayed
   - **Expected:** UI updates (title, meta, and reference in fuzzy finder)

**Test 7: Update entry non-title field**
1. Open Entry List (^P)
2. Select an Entry
3. Update non-title field
   - **Expected:** Field changes
4. Save (^O)
   - **Expected:** Notification displayed
   - **Expected:** UI updates (meta changes)

**Test 8: Update entry group and save**
1. Open Entry List (^P)
2. Select an Entry
3. Update Group
4. Save (^O)
   - **Expected:** Notification displayed
   - **Expected:** UI updates (meta changes)
5. Open Entry List (^P)
   - **Expected:** Updated entry appears in list of entries

**Test 9: Update entry group and cancel**
1. Open Entry List (^P)
2. Select an Entry
3. Update Group
4. Cancel (ESC)
   - **Expected:** Entry returns to previous state
5. Open Entry List (^P)
   - **Expected:** Updated entry does NOT appear in new location

**Test 10: Dismiss entry deletion**
1. Open Entry List (^P)
2. Select an Entry
3. Delete (^D)
4. Say No (N)
5. Open Entry List (^P)
   - **Expected:** Entry still appears in list of entries

**Test 11: Confirm entry deletion**
1. Open Entry List (^P)
2. Select an Entry
3. Delete (^D)
4. Say Yes (Y)
   - **Expected:** Navigates back to list
5. Check Entry List (^P)
   - **Expected:** Entry does NOT appear in list of entries

</details>

<details>
<summary>open: read only, no entry</summary>

**Test: Read-only mode restrictions**
1. Open Entry List (^P)
2. Select Entry
3. Attempt to change fields
   - **Expected:** Cannot change fields (notification displayed)
4. Attempt to change groups
   - **Expected:** Cannot change groups (notification displayed)
5. Attempt to save changes
   - **Expected:** Cannot save changes (notification displayed)
6. Attempt to delete entry
   - **Expected:** Cannot delete entry (notification displayed)

</details>

<details>
<summary>open: no read-only, with ref</summary>

**Test 1: View and copy entry password by reference**
1. Open Entry by Ref
2. Reveal password (^R)
   - **Expected:** Password is readable
3. Copy password (^C)
   - **Expected:** Password is in clipboard

**Test 2: Edit and cancel non-hidden and hidden fields by reference**
1. Open Entry by Ref
2. Update non-hidden field
   - **Expected:** Field changes
3. Cancel (ESC)
   - **Expected:** Entry returns to previous state
4. Update hidden field
   - **Expected:** Error displayed
5. Reveal field (^R)
   - **Expected:** Password can be updated
6. Update password
7. Cancel (ESC)
   - **Expected:** Entry returns to previous state
8. Save (^O)
   - **Expected:** Nothing happens (no changes to save)

**Test 3: Update entry title by reference**
1. Open Entry by Ref
2. Update title field
   - **Expected:** Field changes
3. Save (^O)
   - **Expected:** Notification displayed
   - **Expected:** UI updates (title, meta, and reference in fuzzy finder)

**Test 4: Update entry non-title field by reference**
1. Open Entry by Ref
2. Update non-title field
   - **Expected:** Field changes
3. Save (^O)
   - **Expected:** Notification displayed
   - **Expected:** UI updates (meta changes)

**Test 5: Update entry group and save by reference**
1. Open Entry by Ref
2. Update Group
3. Save (^O)
   - **Expected:** Notification displayed
   - **Expected:** UI updates (meta changes)
4. Open Entry List (^P)
   - **Expected:** Updated entry appears in list of entries

**Test 6: Update entry group and cancel by reference**
1. Open Entry by Ref
2. Update Group
3. Cancel (ESC)
   - **Expected:** Entry returns to previous state
4. Open Entry List (^P)
   - **Expected:** Updated entry does NOT appear in new location

**Test 7: Dismiss entry deletion by reference**
1. Open Entry by Ref
2. Delete (^D)
3. Say No (N)
4. Open Entry List (^P)
   - **Expected:** Entry still appears in list of entries

**Test 8: Confirm entry deletion by reference**
1. Open Entry by Ref
2. Delete (^D)
3. Say Yes (Y)
   - **Expected:** Navigates back to list
4. Check Entry List (^P)
   - **Expected:** Entry does NOT appear in list of entries

</details>

<details>
<summary>list</summary>

**Requirements:**
- Shows all entries in the database
- Can be piped to other commands (e.g., `fzf`)
- Allows copying of output
- Uses command aliases (ls)

**Test: List command functionality**
1. Run list command
   - **Expected:** All entries are displayed
   - **Expected:** Aliases are used in display
2. Pipe output to another command
   - **Expected:** Output can be piped successfully
3. Copy result
4. Paste result
   - **Expected:** Result can be copied and pasted

</details>

<details>
<summary>credentials</summary>

**Requirements:**
- Password can be provided via environment variable
- Archive path can be provided via environment variable
- Password can be typed interactively
- Database can be unlocked using a keyfile

**Test: Credentials handling**
1. Open with password from environment
   - **Expected:** Database opens successfully
2. Open with archive from environment
   - **Expected:** Database opens successfully
3. Open with typed password
   - **Expected:** Database opens successfully
4. Open with keyfile
   - **Expected:** Database opens successfully

</details>

<details>
<summary>copy</summary>

**Requirements:**
- Can copy password field
- Can copy other fields
- Handles non-existing archive with appropriate error

**Test 1: Copy password field**
1. Run copy command for password field
   - **Expected:** Password is copied to clipboard

**Test 2: Copy non-password field**
1. Run copy command for another field
   - **Expected:** Field value is copied to clipboard

**Test 3: Copy from non-existing archive**
1. Run copy command with non-existing archive path
   - **Expected:** Error message displayed (same as non-existing archive error)
  
**Test 4: Copy from non-existing field**
- Steps:
  1. Run copy command with non-existing field
- Expected:
  - Error: `Missing field "{FIELD}" in entry {ENTRY}`

</details>

<details>
<summary>errors</summary>

**Test 1: Non-existing archive**
1. Attempt to open non-existing archive
   - **Expected:** Error: `"open {FILE}: no such file or directory"`

**Test 2: Invalid credentials**
1. Attempt to open archive with wrong password
   - **Expected:** Error: `"Wrong password? HMAC-SHA256 of header mismatching"`

**Test 3: Missing keyfile**
1. Attempt to open archive with non-existing keyfile
   - **Expected:** Error: `"open {FILE}: no such file or directory"`

**Test 4: Invalid keyfile**
1. Attempt to open archive with invalid keyfile
   - **Expected:** Error: `"Wrong password? HMAC-SHA256 of header mismatching"`

</details>
