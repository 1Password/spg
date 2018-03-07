# AgileBits' Strong Password Generator

This is a *work in progress*, and all **interfaces are subject to change**.

The Strong Password Generator package offers the underlying engine for flexible specification of generated password requirements and ensuring that the generated passwords it returns follow a uniform distribution.

The clients of this package are expected to manage what is presented to users. This engine offers far greater flexibility than should normally be exposed to users.

The various `Generate()` methods return a `Password` object which has a `String()` method and an `Entropy()` method. But for the moment (March 2018), how the `Generate()` methods are called and how the generators are build is still subject to change.