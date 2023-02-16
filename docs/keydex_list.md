## keydex list

Lists all the entries in the database

### Synopsis

Lists all the entries in the database. 

The list of references - in the form of - /database/group/.../entry will be printed on stadout, allowing for piping.
The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the KEYDEX_DATABASE environment variable.
This command can be used in conjuction with tools such like 'fzf' or 'dmenu' to browse the databse and pipe the result to other commands.

See "Examples" for more details.

```
keydex list [file] [flags]
```

### Examples

```
  # List all entries of vault.kdbx database
  keydex list vault.kdbx

  # List entries, browse them with fzf and copy the result to the clipboard
  export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
  export KEYDEX_DATABASE=~/vault.kdbx

  keydex list | fzf | keydex copy
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
  -k, --key string   path to the key file to unlock the database
```

### SEE ALSO

* [keydex](keydex.md)	 - Manage KeePass databases from your terminal.

