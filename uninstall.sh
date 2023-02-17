#!/bin/env sh

if [ "$(id -u)" -ne 0 ]; then 
  echo "Uninstallation needs to be run as super user. Please run 'sudo $0' to proceed."
  exit 1
fi

echo "Removing previous installation..."
rm /usr/local/bin/keydex
echo "Done!"