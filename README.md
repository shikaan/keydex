kpcli
===

Keepass compatible CLI

## Usage


```sh
# opens the fuzzy search and copies the password
kpcli browse copy ~/.local/share/vaults/example.kdbx

# opens the fuzzy search and edits the selected entry
kpcli browse edit ~/.local/share/vaults/example.kdbx

# copies password from the referenced entry (or stdinput)
kpcli copy ~/.local/share/vaults/example.kdbx REF

# edits the referenced entry (or stdinput)
kpcli edit ~/.local/share/vaults/example.kdbx REF

# lists all the entries, to be used with fzf
kpcli list ~/.local/share/vaults/example.kdbx
```

Examples

With fzf

```
kpcli list ~/.local/share/vaults/example.kdbx | fzf | kpcli copy ~/.local/share/vaults/example.kdbx  
```

Get entry
```

```
