## keydex

Manage KeePass databases from your terminal.

### Synopsis

keydex is a command line utility to manage KeePass databases. It comes with
a simple, display-oriented editor inspired by the minimalism of GNU nano.

keydex reads the following environment variables:

  - KEYDEX_PASSPHRASE
    When this variable is set, keydex will skip the password prompt. It can
    be replaced by utils such as 'autoexpect'.

  - KEYDEX_DATABASE
    Path to the *.kbdx database to unlock. Providing 'file' inline overrides
    this value.

  - KEYDEX_KEY
    Path to the optional *.key file used to unlock the database. The '--key'
    flag overrides this value.

All the entries are identified by a path-like reference like
/database/group1/../groupN/entry where 'database' is the database name,
'groupN' are the (nested) groups names, and 'entry' is the entry title.

Internally all the entries are referenced by a UUID, however keydex will read
the first occurrence of a reference in cases of conflicts. Writes are always
done via UUID and they are therefore conflict-safe.

Some commands use the system clipboard, in absence of which keydex will fail.

```
keydex [flags]
```

### Options

```
  -h, --help   help for keydex
```

### SEE ALSO

* [keydex copy](keydex_copy.md)	 - Copies a field of a reference to the clipboard.
* [keydex create](keydex_create.md)	 - Create an empty KeePass archive.
* [keydex list](keydex_list.md)	 - Lists all the entries in the database
* [keydex open](keydex_open.md)	 - Open the entry editor for a reference.

