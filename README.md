# spg - A go package for strong password generation

[![GoDoc: Reference](https://godoc.org/go.1password.io/spg?status.svg)](https://godoc.org/go.1password.io/spg) [![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

1Password's Strong Password Generator package offers the underlying engine for flexible specification of generated password requirements and ensuring that the generated passwords it returns follow a uniform distribution.

The clients of this package are expected to manage what is presented to users. This engine offers far greater flexibility than should normally be exposed to users.

## Get started

Use `go get`:

```bash
go get go.1password.io/spg
```



## Vendored dependencies

Before you can successfully build, you may need to install dependencies. These are currently[^1] managed using [`govendor`](https://github.com/kardianos/govendor). Install it if needed,

```
go get -u github.com/kardianos/govendor
```

And then use 

```
govendor sync
```
to fetch the appropriate dependencies into `./vendor`

[^1]: We will probably switch to go modules at some point

## License

1Password's spg is copyright 2018, AgileBits Inc and licensed under [version 2.0 of the Apache License Agreement](./LICENSE).


## Contributing

This is on Github: https://github.com/1password/spg create issues, forks, etc there.