Contributing
---

Practical guidelines to contribute to this project.

## Building from source

In order to build the app from source, you can just type:

```sh
# Get all the dependencies
go get ./...

# Generate app metadata (e.g., revision, version from git history)
make info

# Generate the binary for the current architecture and os
make build
```

## Run locally

The easiest way is to use something like `direnv` and define a `.envrc` file like

```sh
export KEYDEX_PASSPHRASE="my-password"
export KEYDEX_DATABASE="test.kdbx"
```

The `.envrc` and the `test.kdbx` files are gitignored already.

## Sharing a testing build

If you need to share a testing build, maybe because you are working on a complicated
bug, you can simply add your branch name to the [configuration](./.github/workflows/build.yml)
file.

It will create a new Github Release tagged with the branch name. 
