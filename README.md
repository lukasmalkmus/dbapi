# LukasMa/dbapi
> Access the Deutsche Bank API from your go application. - by **[Lukas Malkmus](https://github.com/LukasMa)**

[![Build Status][build_badge]][build]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![GoDoc][docs_badge]][docs]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

---

## Table of Contents
1. [Introduction](#introduction)
2. [Features](#features)
3. [Usage](#usage)
4. [Contributing](#contributing)
5. [License](#license)

### Introduction
This package is a small wrapper around the recently released Deutsche Bank API
and covers all of the available http endpoints.

The Deutsche Bank API will grow over time and so will this package.

### Features
  - [x] Cover all available endpoints
    - [x] Accounts (`/cashAccounts`)
    - [x] Addresses (`/addresses`)
    - [x] Transactions (`/transactions`)
    - [x] UserInfo (`/userInfo`)
  - [x] Select API version to use (currently only v1 is available)
  - [x] Basic test suit
  - [x] Easy to use package API

#### Todo
  - [ ] Provide authentication?

### Usage
#### Requirements
Create an account on the [Developer Portal](https://developer.db.com) and follow
the instructions there. The common workflow is to create an application and at
least one test user to get started.

#### Installation
Please use a dependency manager like [glide](http://glide.sh) to make sure you
use a tagged release.

Install using `go get`:
```bash
go get -u github.com/LukasMa/dbapi
```

#### Usage
##### Authentication
The Deutsche Bank API is secured by OAuth2 and you need an access token to
retrieve data from the endpoints. This package **does not** provide an OAuth2
client since there are many good implementations out there. To use the `dbapi`
client you need the `Access Token`.

##### Creating a new api client.
To retrieve data you need to create a new client:
```go
import github.com/LukasMa/dbapi

const AccessToken = "..."

api, err := dbapi.New(
    dbapi.SetToken(AccessToken),
)
if err != nil {
    log.Fatalln(err)
}
```

Since the access token is bound to a specific user the client can only scrape
data from exactly this user.

_Please note that the user must grant the correct rights (`scopes`) during the
authentication process or you might not be allowed to access the corresponding
api endpoints._

**It is also possible to use a custom http client to make requests:**
```go
// Create your custom http client.
client := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        Dial: (&net.Dialer{
            Timeout:   3 * time.Second,
            KeepAlive: 30 * time.Second,
        }).Dial,
        ExpectContinueTimeout: 1 * time.Second,
        ResponseHeaderTimeout: 3 * time.Second,
        TLSHandshakeTimeout:   3 * time.Second,
    },
}

// Use your custom http client.
api, err := dbapi.New(
    dbapi.SetToken(AccessToken),
    dbapi.SetClient(client),
)
// ...
```

**Options can also be applied to a client instance if it has already been created:**
```go
api, err := dbapi.New()

api.Options(
    dbapi.SetToken(AccessToken),
)
```

##### Accessing resources
Accessing the endpoints is easy. Since the API is in an early state there aren't
many enpoints, yet. A list of available endpoints can be found on the
Developer Portal > API Explorer. Or take a look at the [swagger specification](https://simulator-api.db.com/gw/dbapi/v1/swagger.json).

```go
accounts, response, err := api.Accounts.GetAll()
if err != nil {
    fmt.Println(response)
    log.Fatalln(err)
}
fmt.Printf("%v", accounts)
```

There are currently four endpoints with either one or two methodes to explore:
  - `api.Addresses`
    - `.Get()`
  - `api.Accounts`
    - `.Get(iban string)`
    - `.GetAll()`
  - `api.Transactions`
    - `.Get(iban string)`
    - `.GetAll()`
  - `api.UserInfo`
    - `.Get()`

### Contributing
Feel free to submit PRs or to fill Issues. Every kind of help is appreciated.

### License
(c) Lukas Malkmus, 2017

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.


[docs]: https://godoc.org/github.com/LukasMa/dbapi
[docs_badge]: https://godoc.org/github.com/LukasMa/dbapi?status.svg
[release]: https://github.com/LukasMa/dbapi/releases
[release_badge]: https://img.shields.io/github/release/LukasMa/dbapi.svg
[report]: https://goreportcard.com/report/github.com/LukasMa/dbapi
[report_badge]: https://goreportcard.com/badge/github.com/LukasMa/dbapi
[build]: https://travis-ci.org/LukasMa/dbapi
[build_badge]: https://travis-ci.org/LukasMa/dbapi.svg
[coverage]: https://coveralls.io/github/LukasMa/dbapi?branch=master
[coverage_badge]: https://coveralls.io/repos/github/LukasMa/dbapi/badge.svg?branch=master
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/badge/license-MIT-blue.svg
