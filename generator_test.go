//
// generator_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package enumer

import "testing"

func TestGenerator(t *testing.T) {
	g := Generator{}
	g.parse("testdata/day")
	g.generate("Day")
}
