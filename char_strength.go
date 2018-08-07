package spg

import (
	"math/big"

	set "github.com/deckarep/golang-set"
)

// This is where we do the math for the entropy calculation for character
// passwords. The trick is for when a character is _required_ from a particular set

func (r CharRecipe) entropyWithRequired() float32 {

	// place holder until other parts written
	return 0.0
}

// N is the number of possible passwords that can be generated.
// Unfortunately, we can't take the log until the very end, so we will
// be dealing with some very large numbers.
func N(aSet set.Set, rSet set.Set, n int) *big.Int {
	out := &big.Int{}

	// stuff will go here

	return out
}
