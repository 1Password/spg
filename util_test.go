package spg

import (
	"math"
	"testing"
)

func TestRandomUint32(t *testing.T) {
	// t.Errorf("%d", uint32(-5))
	c := 0
	for i := 1; i <= 20; i++ {
		n := randomUint32()
		if n > math.MaxInt32 {
			c++
			t.Errorf("Generated big number %d: %d", c, n)
		}
	}
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
