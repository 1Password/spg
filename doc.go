/*
Package spg provides 1Password's Strong Password Generator engine for uniform and flexible password generation.


The Strong Password Generator package offers the underlying engine for flexible
specification of generated password requirements and ensuring that the generated
passwords it returns follow a uniform distribution.

The clients of this package are expected to manage what is presented to users.
This engine offers far greater flexibility than should normally be exposed to users.

Usage overview

The very highlevel usage over view for most cases is

1. The user will create either a character recipe, r, with NewCharRecipe()
or a word (or syllable) list recipe, r, with NewWLRecipe().

2. The user will call the r.Generate() method of a recipe, r, to generate a password, pwd.

3. The returned password, pwd, has a String() method, which does the obvious thing and
Entropy field, which contains the min-Entropy based on the recipe.

Wordlist and pronounceable

The word list generator produces things like "correct horse battery staple", but
when the list is of pronounceable syllables, it can also be set up to produce things
like

    Mirk9vust8jilk3rooy
    scuy9lam2lerk9Kais
    smoh1fock6mirn7Lic
    jaud3Rew4jo6mont

Lengths for these are specified in terms of the number of elements drawn from the
list to be included in these passwords (not counting the separators).
Although the above examples all have different lengths in terms of number of characters,
they were all specified as Length 4.

The passwords that one gets depend on the word list recipe, WLRecipe, and the actual
word list provided.

Character passwords

Character-based are your typical notion of generated password,
however these can be specified in ways to produce only numeric PINs if desired.
The passwords generated are a function of the CharRecipe.

The Generate and Entropy methods

The word list and character recipes (WLRecipe, CharRecipe) implement a Generator
interface with two methods, Generate and Entropy.

Generate returns a Password. There is a fair amount of internal structure
to a Password object, but the ones you are most after is available through
the Password.String() method and the Entropy field.

Entropy returns the entropy of a password that would be generated
given the current recipe.


A word about Entropy

Entropy is a highly misleading concept when applied to passwords. In the general case it
is either an incoherent concept or the wrong concept to use when talking about the strength
of a password.
It does, however, make sense when a password is drawn uniformly from a space of possible passwords.
This package does ensure that passwords are generated uniformly given the recipe
passed to the generator, with the exception of the interaction of capitalizaton for some wordlists.
In those cases, min-entropy is reported.

Entropy is a function solely of the recipe.
*/
package spg

// This file is for package documentation only
