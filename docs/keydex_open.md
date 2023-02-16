## keydex open

Open the entry editor for a reference.

### Synopsis

Open the entry editor for a reference.

Reads a 'reference' from the database at 'file' and opens the editor there. If no reference is passed, it opens a fuzzy search within the editor.

The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the KEYDEX_DATABASE environment variable.
The 'reference' can be passed as last argument; if the reference is missing, it opens a fuzzy search.
Use the 'list' command to get a list of all the references in the database.

See "Examples" for more details.

```
keydex open [file] [reference] [flags]
```

### Examples

```
  # Opens the "github" entry in the "coding" group in the "test" database at test.kdbx
  keydex open test.kdbx /test/coding/github
  
  # Open fuzzy search within the test.kdbx database
  keydex open test.kdbx

  # Or with environment variables
  export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
  export KEYDEX_DATABASE=test.kdbx
  keydex open

  # List entries, browse them with fzf and edit the result
  export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
  export KEYDEX_DATABASE=test.kdbx

  keydex list | fzf | keydex open
```

### Options

```
  -h, --help   help for open
```

### Options inherited from parent commands

```
  -k, --key string   path to the key file to unlock the database
```

### SEE ALSO

* [keydex](keydex.md)	 - Manage KeePass databases from your terminal.

