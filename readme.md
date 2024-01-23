# `ireturn` [![Stand with Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://u24.gov.ua/) [![Code Coverage](https://coveralls.io/repos/github/butuzov/ireturn/badge.svg?branch=main)](https://coveralls.io/github/butuzov/ireturn?branch=main) [![build status](https://github.com/butuzov/ireturn/actions/workflows/main.yaml/badge.svg?branch=main)]() [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

Accept Interfaces, Return Concrete Types

## Install

You can get `ireturn` with `go install` command. Go1.18+ required.

```shell
go install github.com/butuzov/ireturn/cmd/ireturn@latest
```

### Compiled Binary

Or you can download the suitable binary from the [releases](https://github.com/butuzov/ireturn/releases) section.

## Usage

`ireturn` work with two arguments (but allow to use of only one of them in the same moment):

* `accept` - accept-list of the comma-separated interfaces.
* `reject` - reject-list of the comma-separated interfaces.

By default, `ireturn` will accept all errors (`error`), empty interfaces (`interfaces{}`), anonymous interfaces declarations ( `interface { methodName() }` ) and interfaces from standard library as a valid ones.

Interfaces in the list can be provided as regexps or keywords ("error" for "error", "empty" for `interface{}`, `anon` for anonymous interfaces):

```bash
# allow usage of empty interfaces, errors and Doer interface from any package.
ireturn --accept="\\.Doer,error,empty" ./...
# reject standard library interfaces and plinko.Payload as valid ones
ireturn --reject="std,github.com/shipt/plinko.Payload" ./...
# default settings allows errors, empty interfaces, anonymous declarations and standard library
ireturn ./...
# checkfor non idiomatic interface names
ireturn -allow="error,generic,anon,stdlib,.*(or|er)$" ./...
```

### Keywords

You can use shorthand for some types of interfaces:

* `empty` for `interface{}` type
* `anon` for anonymous declarations `interface{ someMethod() }`
* `error` for `error` type
* `stdlib` for all interfaces from standard library.
* `generic` for generic interfaces (added in go1.18)

### Disable directive

`golangci-lint` compliant disable directive `//nolint: ireturn` can be used with `ireturn`

### GitHub Action

```
- uses: butuzov/ireturn-linter@main
  with:
    allow: "error,empty"
```

## Examples

```go
// Bad.
type Doer interface { Do() }
type IDoer struct{}
func New() Doer { return new(IDoer)}
func (d *IDoer) Do() {/*...*/}

// Good.
type Doer interface { Do() }
type IDoer struct{}
func New() *IDoer { return new(IDoer)}
func (d *IDoer) Do() {/*...*/}

// Very Good (Verify Interface Compliance in compile time)
var _ Doer = (*IDoer)(nil)

type Doer interface { Do() }
type IDoer struct{}
func New() *IDoer { return new(IDoer)}
func (d *IDoer) Do() {/*...*/}

```

## Reading List
* [Rob Pike's comment on "Accept Interfaces Return Struct in Go"](https://github.com/go-proverbs/go-proverbs.github.io/issues/37)
* [Accept Interfaces Return Struct in Go](https://mycodesmells.com/post/accept-interfaces-return-struct-in-go)
* [Accept Interface Return Struct](https://blog.dlow.me/programming/golang/accept-interface-return-struct/)
* [What “accept interfaces, return structs” means in Go](https://medium.com/@cep21/what-accept-interfaces-return-structs-means-in-go-2fe879e25ee8)
