#!/bin/bash

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$([[ $(uname -m) == "x86_64" ]] && echo "amd64" || echo "386")

has_wget=$(which wget)

echo "Installing the binary..."
if [ -z has_wget ]; then
  curl -sfLo /usr/local/bin/keydex https://github.com/shikaan/keydex/releases/latest/download/keydex-${OS}-${ARCH}
else
  wget -q -O /usr/local/bin/keydex https://github.com/shikaan/keydex/releases/latest/download/keydex-${OS}-${ARCH}
fi

chmod u+x /usr/local/bin/keydex
chown $(logname) /usr/local/bin/keydex

echo "Done!"