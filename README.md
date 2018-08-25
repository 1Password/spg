# AgileBits' Strong Password Generator

This is a *work in progress*, and all **interfaces are subject to change**.

The Strong Password Generator package offers the underlying engine for flexible specification of generated password requirements and ensuring that the generated passwords it returns follow a uniform distribution.

The clients of this package are expected to manage what is presented to users. This engine offers far greater flexibility than should normally be exposed to users.

## Vendored packages

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