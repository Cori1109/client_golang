// Copyright (c) 2012, Matt T. Proud
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// utility_test.go provides a test suite for all tests in the utility package
// hierarchy.  It employs the gocheck framework for test scaffolding.

package utility

import (
	. "launchpad.net/gocheck"
	"testing"
)

type S struct{}

var _ = Suite(&S{})

func TestUtility(t *testing.T) {
	TestingT(t)
}
