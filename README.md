# Go CSV

CSV encoder written in and for the Go programming language.

## TODOs

- [ ] set up CI
    - [ ] build
    - [ ] code quality - golangci-lint
    - [ ] correctness - tests
- [ ] set up CD
    - [ ] push to go.dev on tag on `main`
- [ ] create default structs of top level functions (similar to `DefaultHTTPClient`):
    - [ ] default reader
    - [ ] default parser
- [ ] rename the variables in [unmarshal](unmarshal.go)
- [ ] make error types part of API
    - [ ] add tests
    - [ ] use sentinel errors or custom error structs
- [ ] parse other value types
    - [ ] handle custom types, either:
        - [ ] 1. check for interface implmentation
        - [ ] 2. register type parsers

## Future Ideas

- [ ] unmarshal single structs
- [ ] unmarshal [][]string (raw?)
- [ ] find a way to attach raw header and string to struct / interface (?)

