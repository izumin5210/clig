# clig / clib
[![Build Status](https://travis-ci.com/izumin5210/clig.svg?branch=master)](https://travis-ci.com/izumin5210/clig)
[![GoDoc](https://godoc.org/github.com/izumin5210/clig/pkg/cli?status.svg)](https://godoc.org/github.com/izumin5210/clig/pkg/cli)
[![license](https://img.shields.io/github/license/izumin5210/clig.svg)](./LICENSE)

- :building_construction: `clig` is a boilerplate generator for CLI tools in Go
- :wrench: [`pkg/clib`](https://godoc.org/github.com/izumin5210/clig/pkg/clib) is a utility package to create CLI tools efficiently 

## Usage

```console
# initialize a new project
$ clig init awesomecli

$ cd awesomecli
$ tree -I 'bin|vendor'
.
├── Gopkg.lock
├── Gopkg.toml
├── Makefile
├── cmd
│   └── awesomecli
│       └── main.go
├── pkg
│   └── awesomecli
│       ├── cmd
│       │   └── cmd.go
│       ├── config.go
│       └── context.go
└── tools.go
```

## Installation
To install clig, you can use `go get`:

```console
$ go get github.com/izumin5210/clig/cmd/clig
```

## Author
- Masayuki Izumi ([@izumin5210](https://github.com/izumin5210))


## License
licensed under the MIT License. See [LICENSE](./LICENSE)
