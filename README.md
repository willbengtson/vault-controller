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

## Examples

### Example Username & Password Auth Backend

1. Stand up a local [Vault](https://www.vaultproject.io/) server listening on default port - 127.0.0.1:8200
2. Enable the [`userpass`](https://www.vaultproject.io/docs/auth/userpass.html) Auth backend
3. Create a user called `testuser` with a password of `testpassword` and assign the default policy
4. Run `./login.sh`
5. Run `./vault-controller --addr=http://127.0.0.1:8200 --log_level=info --log_fmt=json --token_file=.vault-response` to enable renewal of your Vault client token

### Example AWS-EC2

1. Stand up Vault with [`AWS EC2`](https://www.vaultproject.io/docs/auth/aws-ec2.html) Auth Backend
2. Authorized a role called `example prod` with a 12 hour ttl
3. Launch an EC2 instance and authenticate with Vault.  Write the output to `/vault/.vault-response`
4. Install vault-controller as a service
4. Start vault-controller - `service vault-controller start` to auto renew the Vault client token

### Example Vault authentication response
```
{"request_id":"1b105e2e-4996-77e6-a4b7-37583b124ea2","lease_id":"","renewable":false,"lease_duration":0,"data":null,"wrap_info":null,"warnings":null,"auth":{"client_token":"2c89f2ec-d50f-4821-9032-de2664589d26","accessor":"5cfa295b-2b5c-44bd-86c7-67ce90741e37","policies":["default","service/gitlab-dev"],"metadata":{"account_id":"123456789123","ami_id":"ami-00000000","instance_id":"i-837dadd923847109a","nonce":"Gtp5msp3jEq759qpIvntP","region":"us-west-2","role":"example-prod","role_tag_max_ttl":"0s"},"lease_duration":43200,"renewable":true}}
```

### Example vault-controller output

```
{"level":"info","msg":"Starting vault-controller...","time":"2017-03-09T00:01:49Z"}
{"level":"info","msg":"Reading vault secret file from /vault/.vault-response","time":"2017-03-09T00:01:49Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-09T00:01:50Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-09T06:01:52Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-09T12:01:54Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-09T18:01:57Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-10T00:01:59Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-10T06:02:01Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-10T12:02:03Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-10T18:02:05Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-11T00:02:06Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-11T06:02:08Z"}
{"level":"info","msg":"token-renew: Successfully renewed the client token; next renewal in 21600 seconds","time":"2017-03-11T12:02:10Z"}
```