## keydex diff

Compares two KeePass archives

### Synopsis

Compares two KeePass archives and outputs which entries were added, removed,
or modified. Output follows the unified diff format so it can be piped into other tools.

The 'file-a' and 'file-b' arguments are paths to the *.kdbx archives to compare.
Passphrases can be provided via environment variables to avoid interactive prompts.

```
keydex diff [file-a] [file-b] [flags]
```

### Examples

```
  # Compare two archives
  keydex diff old.kdbx new.kdbx

  # Or with environment variables
  export KEYDEX_PASSPHRASE_A=${PASSPHRASE_A}
  export KEYDEX_PASSPHRASE_B=${PASSPHRASE_B}
  keydex diff old.kdbx new.kdbx

  # With key files
  keydex diff --key-a old.key --key-b new.key old.kdbx new.kdbx
```

### Options

```
  -h, --help           help for diff
      --key-a string   path to the key file for the first archive
      --key-b string   path to the key file for the second archive
```

### SEE ALSO

* [keydex](keydex.md)	 - Manage KeePass databases from your terminal.

