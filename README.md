kpcli
===

Keepass compatible CLI

## Usage


```sh
# opens the editor at the given location
kpcli open ~/.local/share/vaults/example.kdbx /example/group/entry

# copies password from the referenced entry (or stdinput)
kpcli copy-password ~/.local/share/vaults/example.kdbx /example/group/entry

# lists all the entries, to be used with fzf
kpcli list ~/.local/share/vaults/example.kdbx /example/group/entry 
```

## Examples

With fzf

```
# Copy entry's password selected with fzf in the clipboard
kpcli list ~/.local/share/vaults/example.kdbx | fzf | kpcli copy ~/.local/share/vaults/example.kdbx  

# Open entry at ref selected with fzf
kpcli list ~/.local/share/vaults/example.kdbx | fzf | kpcli open ~/.local/share/vaults/example.kdbx  
```

Get entry
```

```
