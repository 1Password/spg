# 1Password's Strong Password Generator

The Strong Password Generator package offers the underlying engine for flexible specification of generated password requirements and ensuring that the generated passwords it returns follow a uniform distribution.

The clients of this package are expected to manage what is presented to users. This engine offers far greater flexibility than should normally be exposed to users.

Use `go doc` for package documentation, or `godoc -http=:6060` to run a documentation server on http://localhost:6060

## License

1Password's spg is copyright 2018, AgileBits Inc and licensed under [version 2.0 of the Apache License Agreement](./LICENSE).

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