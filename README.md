# Go CSV

[![Build](https://github.com/JoelLau/go-csv/actions/workflows/build.yml/badge.svg)](https://github.com/JoelLau/go-csv/actions/workflows/build.yml)
[![Test](https://github.com/JoelLau/go-csv/actions/workflows/test.yml/badge.svg)](https://github.com/JoelLau/go-csv/actions/workflows/test.yml)

CSV encoder written in and for the Go programming language.

## TODOs

- [-] set up CI
    - [x] build
    - [ ] code quality - golangci-lint
    - [x] correctness - tests
- [x] ~~set up CD~~ done automatically
    - [x] ~~push to go.dev on tag on `main`~~
- [ ] create default structs of top level functions (similar to `DefaultHTTPClient`):
    - [ ] default reader
    - [ ] default parser
- [ ] rename the variables in [unmarshal](unmarshal.go)
- [ ] make error types part of API
    - [ ] add tests
    - [ ] use sentinel errors or custom error structs
- [x] parse other value types
    - [x] basic types
    - [x] unmapped columns
    - [ ] pointers (for nullsy values)
    - [ ] handle custom types, either:
        - [ ] 1. check for interface implmentation
        - [ ] 2. register type parsers

## Future Ideas

- [ ] unmarshal single structs
- [ ] unmarshal [][]string (raw?)
- [ ] find a way to attach raw header and string to struct / interface (?)

