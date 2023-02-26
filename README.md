<p align="center">
  <img width="96" height="96" src="./docs/96x96.png" alt="logo">
</p>

<h1 align="center">keydex</h1>

<p align="center">
Terminal based password manager for KeePass databases.
</p>

## ⚡️ Quick start

### Installation

_MacOS and Linux_
```sh
sudo sh -c "curl -s https://shikaan.github.io/sup/install | REPO=shikaan/keydex sh -"

# or

sudo sh -c "wget -q https://shikaan.github.io/sup/install -O- | REPO=shikaan/keydex sh -"
```

_Windows and manual instructions_

Head to the [releases](https://github.com/shikaan/keydex/releases) page and download the executable for your system and architecture.

### Usage

You can get started simply [opening](./docs/keydex_open.md) your database.

```sh
# opens the interactive editor
keydex open ~/example.kdbx
```

However, the most common use case for `keydex` is [copying](./docs/keydex_copy.md) a password to your clipboard.

```sh
# copies password from the referenced entry (or stdinput)
keydex copy ~/example.kdbx /example/group/entry
```

Using environment variables and aliases you can save a couple of keystrokes

```sh
# ~/.bashrc or ~/.zshrc
export KEYDEX_PASSPHRASE=${MY_SECRET_PHRASE}
export KEYDEX_DATABASE=~/example.kdbx

alias entry-pwd="keydex copy /example/group/entry"

# and then you can use it like
$ entry-pwd
```

Opening a given entry and [listing](./docs/keydex_list.md) accept environment variables too.

```sh
# opens the editor at the given location
keydex open /example/group/entry

# lists all the entries in the database
keydex list
```

### Interoperability

keydex was designed to integrate in your existing workflow: it accepts inputs from stdin and can be piped to your existing toolchain. 

For example, here's an of how you can use it to browse entries with [fzf](https://github.com/junegunn/fzf)

```sh
# copy entry's password selected with fzf to the clipboard
keydex list | fzf | keydex copy  

# open entry at ref selected with fzf
keydex list | fzf | keydex open  
```

## 📄 Documentation

More detailed documentation can be found [here](./docs/keydex.md).


## 🤓 Contributing

Have a look through existing [Issues](https://github.com/shikaan/keydex/issues) and [Pull Requests](https://github.com/shikaan/keydex/pulls) that you could help with. If you'd like to request a feature or report a bug, please create a [GitHub Issue](https://github.com/shikaan/keydex/issues).

## License

[MIT](./LICENSE)
