## keydex copy

Copies a field of a reference to the clipboard.

### Synopsis

Copies a field of a reference to the clipboard.

Reads a 'reference' from the database at 'file' and copies the value of 'field' to the clipboard.

The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the KEYDEX_DATABASE environment variable.
The 'reference' can be passed either as the last argument, or can be read from stdin - to allow piping.
Use the 'list' command to get a list of all the references in the database.

See "Examples" for more details.

```
keydex copy [file] [reference] [flags]
```

### Examples

```
  # Copy the password of the "github" entry in the "coding" group in the "test" database at test.kdbx
  keydex copy test.kdbx /test/coding/github
  
  # Or copy the username instead
  keydex copy -f username test.kdbx /test/coding/github

  # Or with stdin
  export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
  echo "/test/coding/github" | keydex copy test.kdbx

  # Or with stdin and environment variables
  export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
  export KEYDEX_DATABASE=test.kdbx
  echo "/test/coding/github" | keydex copy

  # List entries, browse them with fzf and copy the username to the clipboard
  export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
  export KEYDEX_DATABASE=test.kdbx

  keydex list | fzf | keydex copy -f username
```

### Options

```
  -f, --field string   field whose value will be copied (default "password")
  -h, --help           help for copy
```

### Options inherited from parent commands

```
  -k, --key string   path to the key file to unlock the database
```

### SEE ALSO

* [keydex](keydex.md)	 - Manage KeePass databases from your terminal.

