## keydex create

Create an empty KeePass archive.

### Synopsis

Create an empty KeePass archive.

Creates a new KeePass database at 'file' called 'name'. You will be prompted to set a passphrase for the new database.

See "Examples" for more details.

```
keydex create [file] [name] [flags]
```

### Examples

```
  # Create a new database called "vault" at vault.kdbx
  keydex create vault.kdbx vault

  # Create a new database at a specific path
  keydex create ~/passwords/work.kdbx work
```

### Options

```
  -h, --help   help for create
```

### SEE ALSO

* [keydex](keydex.md)	 - Manage KeePass databases from your terminal.

