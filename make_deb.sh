#!/bin/sh

VERSION="0.0.1"
CONTACT="Jane Doe <janed@example.com>"
PACKAGE_NAME="vault-controller"

DIRNAME="$(cd "$(dirname "$0")" && pwd)"
OLDESTPWD="$PWD"

go build
rm -rf "$PWD/rootfs"
mkdir -p "$PWD/rootfs/usr/local/bin"
mv "$PWD/vault-controller" "$PWD/rootfs/usr/local/bin/"

fakeroot fpm -C "$PWD/rootfs" \
    --license "MIT" \
    --url "https://github.com/willbengtson/vault-controller" \
    --vendor "" \
    --description "vault-controller is a binary that manages renewing vault tokens." \
    -m "${CONTACT}" \
    -n "${PACKAGE_NAME}" -v "$VERSION" \
    -p "$OLDESTPWD/${PACKAGE_NAME}_${VERSION}_linux_amd64.deb" \
    -s "dir" -t "deb" \
    "usr"