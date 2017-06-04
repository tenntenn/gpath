# gpath [![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godoc] [![Travis](https://img.shields.io/travis/tenntenn/gpath.svg?style=flat-square)][travis] [![Go Report Card](https://goreportcard.com/badge/github.com/tenntenn/gpath)](https://goreportcard.com/report/github.com/tenntenn/gpath) [![codecov](https://codecov.io/gh/tenntenn/gpath/branch/master/graph/badge.svg)](https://codecov.io/gh/tenntenn/gpath)

[godoc]: http://godoc.org/github.com/tenntenn/gpath
[travis]: https://travis-ci.org/tenntenn/gpath

`gpath` is a Go package to access a field by a path using `reflect` pacakge.

A path is represented by a Go's expression such as `A.B.C[0]`.
You can use selector and index expressions into a path.

See usage and example in [GoDoc](https://godoc.org/github.com/tenntenn/gpath).

*NOTE*: This package is experimental and may make backward-incompatible changes.

## Install

Use go get:

```
$ go get github.com/tenntenn/gpath
```

## Usage

All usage are described in [GoDoc](https://godoc.org/github.com/tenntenn/gpath).

[mercari/go-httpdoc](https://github.com/mercari/go-httpdoc) is a good example for gpath.
