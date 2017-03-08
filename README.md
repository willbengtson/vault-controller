# Vault Controller

[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT)

## About

vault-controller is a binary built to manage the HashiCorp Vault client token for you.  vault-controller will take the Vault auth response manage your TTL for you.  This is an inspiration from [Kelsey Hightower](https://github.com/kelseyhightower) and his vault-controller for managing renewal in kubernetes pods.

##### Goals

* Easy: Make it easy to use
* Eliminate the need to `echo "0 */5 * * * root /usr/local/bin/vault token-renew" > /etc/cron.d/renew-vault-token`

## Usage

##### Installation

1. Install [golang](https://golang.org/doc/install), version 1.7 or greater is recommended
2. Install [`govendor`](https://github.com/kardianos/govendor) if you haven't already

    ```go get -u github.com/kardianos/govendor```
3. Clone the repo

    ```
    git clone (this repo)
    cd vault-controller
    ```
    
4. Build the binary

    ```
    go build
    ```

5. Copy the binary `vault-controller` to wherever you'd like

##### Building a debian package

If you'd like to build a debian package you need to make sure you have `fpm` installed first:

1. Install [`fpm`](https://github.com/jordansissel/fpm)
2. Run `./make_deb.sh`

## Example

1. Stand up a local [Vault](https://www.vaultproject.io/) server listening on default port - 127.0.0.1:8200
2. Enable the [`userpass`](https://www.vaultproject.io/docs/auth/userpass.html) Auth backend
3. Create a user called `testuser` with a password of `testpassword` and assign the default policy
4. Run `./login.sh`
5. Run `./vault-controller --addr=http://127.0.0.1:8200 --log_level=info --log_fmt=json --token_file=.vault-response` to enable renewal of your Vault client token
