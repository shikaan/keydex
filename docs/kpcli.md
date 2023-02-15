## kpcli

Manage KeePass databases from your terminal.

### Synopsis

kpcli is a simple, display-oriented browser and editor for KeePass databases. The user interface is highly inspired by the minimalism of GNU nano: commands are displayed at the bottom of the screen, and context-sensitive help is provided.

Commands are inserted using control-key (^) combinations. For example, "^C" means "Ctrl+C". kpcli comes with subcommands to read and write entries in the provided database. More information available at "kpcli help [command]". 

To facilitate scripting, this tool comes with the ability of reading the following environment variables:

  - KPCLI_PASSPHRASE 
    When this variable is set, kpcli will skip the password prompt. It can be replaced by utils such as 'autoexpect'.

  - KPCLI_DATABASE
    Is the path to the *.kbdx database to unlock. Providing 'file' inline overrides this value.

  - KPCLI_KEY
    Is the path to the optional *.key file used to unlock the database. Providing the '--key' flag inline overrides this value.

All the entries are referenced with a path-like reference string shaped like /database/group1/../groupN/entry where 'database' is the database name, 'groupX' is the group name, and 'entry' is the entry title. Internally all the entries are referenced by a UUID, however kpcli will read the first occurrence of a reference in cases of conflicts. Writes are always done via UUID and they are threfore conflict-safe.
    
Some commands make use of the system clipboard, in absence of which the command will silently fail.

More specific help is available contextually or by typing "kpcli help [command]".

```
kpcli [flags]
```

### Options

```
  -h, --help         help for kpcli
  -k, --key string   path to the key file to unlock the database
  -v, --version      print the version number of kpcli.
```

### SEE ALSO

* [kpcli copy](kpcli_copy.md)	 - Copies the password of a reference to the clipboard.
* [kpcli list](kpcli_list.md)	 - Lists all the entries in the database
* [kpcli open](kpcli_open.md)	 - Open the entry editor for a reference.
* [kpcli version](kpcli_version.md)	 - Print the version number of kpcli.

