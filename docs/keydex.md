## keydex

Manage KeePass databases from your terminal.

### Synopsis

keydex is a command line utility to manage KeePass databases. It comes with subcommands for managing the entries and a simple, display-oriented editor inspired by the minimalism of GNU nano.

keydex can read the following environment variables:

  - KEYDEX_PASSPHRASE 
    When this variable is set, keydex will skip the password prompt. It can be replaced by utils such as 'autoexpect'.

  - KEYDEX_DATABASE
    Is the path to the *.kbdx database to unlock. Providing 'file' inline overrides this value.

  - KEYDEX_KEY
    Is the path to the optional *.key file used to unlock the database. Providing the '--key' flag inline overrides this value.

All the entries are referenced with a path-like reference string shaped like /database/group1/../groupN/entry where 'database' is the database name, 'groupX' is the group name, and 'entry' is the entry title. 

Internally all the entries are referenced by a UUID, however keydex will read the first occurrence of a reference in cases of conflicts. Writes are always done via UUID and they are threfore conflict-safe.
    
Some commands make use of the system clipboard, in absence of which keydex will fail.

```
keydex [flags]
```

### Options

```
  -h, --help         help for keydex
  -k, --key string   path to the key file to unlock the database
```

### SEE ALSO

* [keydex copy](keydex_copy.md)	 - Copies the password of a reference to the clipboard.
* [keydex list](keydex_list.md)	 - Lists all the entries in the database
* [keydex open](keydex_open.md)	 - Open the entry editor for a reference.

