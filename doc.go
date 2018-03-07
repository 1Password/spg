/*
Package spg provides AgileBits' Strong Password Generator engine

The Strong Password Generator package offers the underlying engine for flexible
specification of generated password requirements and ensuring that the generated
passwords it returns follow a uniform distribution.

The clients of this package are expected to manage what is presented to users.
This engine offers far greater flexibility than should normally be exposed to users.

The various `Generate()` methods return a `Password` object which has a `String()` method
and an `Entropy()` method.

Work in Progress -- Interface may change

For the moment (March 2018), how the `Generate()` methods are called and how the
generators are build is still subject to change.

A word about "Entropy"

Entropy is a highly misleading concept when applied to passwords. In the general case it
is either an incoherent concept or the wrong concept to use when talking about the strength
of a password.
It does, however, make sense when a password is drawn uniformly from a space of possible passwords.
This package does ensure that passwords are generated uniformly given the requirements
passed to the generator. Indeed, the Entropy is a function solely of those requirements
and paramater.
*/
package spg

// This file is for package documentation only
