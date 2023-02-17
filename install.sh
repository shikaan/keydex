#!/bin/env sh

get_architecture () {
  arch="$(uname -m)"

  if [ "$arch" = "x86_64" ]; then
    echo "am64"
  else
    echo "$arch"
  fi
}

get_os () {
  uname | tr '[:upper:]' '[:lower:]'
}

OS=$(get_os)
ARCH=$(get_architecture)
HAS_WGET=$(which wget)

if [ "$(id -u)" -ne 0 ]; then 
  echo "Installation needs to be run as super user. Please run 'sudo $0' to proceed."
  exit 1
fi

echo "> Downloading..."
if [ -z "$HAS_WGET" ]; then
  curl -sfLo /usr/local/bin/keydex https://github.com/shikaan/keydex/releases/latest/download/keydex-"${OS}"-"${ARCH}"
else
  wget -q -O /usr/local/bin/keydex https://github.com/shikaan/keydex/releases/latest/download/keydex-"${OS}"-"${ARCH}"
fi

echo "> Setting permissions..."
chmod u+x /usr/local/bin/keydex
chown "$(logname)" /usr/local/bin/keydex

echo "Done!"