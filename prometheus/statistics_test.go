// Copyright (c) 2013, Prometheus Team
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	. "github.com/matttproud/gocheck"
)

func (s *S) TestAverageOnEmpty(c *C) {
	empty := []float64{}
	var v float64 = AverageReducer(empty)

	c.Assert(v, IsNaN)
}

func (s *S) TestAverageForSingleton(c *C) {
	input := []float64{5}
	var v float64 = AverageReducer(input)

	c.Check(v, Equals, 5.0)
}

func (s *S) TestAverage(c *C) {
	input := []float64{5, 15}
	var v float64 = AverageReducer(input)

	c.Check(v, Equals, 10.0)
}

func (s *S) TestFirstModeOnEmpty(c *C) {
	input := []float64{}
	var v float64 = FirstModeReducer(input)

	c.Assert(v, IsNaN)
}

func (s *S) TestFirstModeForSingleton(c *C) {
	input := []float64{5}
	var v float64 = FirstModeReducer(input)

	c.Check(v, Equals, 5.0)
}

func (s *S) TestFirstModeForUnimodal(c *C) {
	input := []float64{1, 2, 3, 4, 3}
	var v float64 = FirstModeReducer(input)

	c.Check(v, Equals, 3.0)
}

func (s *S) TestNearestRankForEmpty(c *C) {
	input := []float64{}

	c.Assert(nearestRank(input, 0), IsNaN)
	c.Assert(nearestRank(input, 50), IsNaN)
	c.Assert(nearestRank(input, 100), IsNaN)
}

func (s *S) TestNearestRankForSingleton(c *C) {
	input := []float64{5}

	c.Check(nearestRank(input, 0), Equals, 5.0)
	c.Check(nearestRank(input, 50), Equals, 5.0)
	c.Check(nearestRank(input, 100), Equals, 5.0)
}

func (s *S) TestNearestRankForDouble(c *C) {
	input := []float64{5, 5}

	c.Check(nearestRank(input, 0), Equals, 5.0)
	c.Check(nearestRank(input, 50), Equals, 5.0)
	c.Check(nearestRank(input, 100), Equals, 5.0)
}

func (s *S) TestNearestRankFor100(c *C) {
	input := make([]float64, 100)

	for i := 0; i < 100; i++ {
		input[i] = float64(i + 1)
	}

	c.Check(nearestRank(input, 0), Equals, 1.0)
	c.Check(nearestRank(input, 50), Equals, 51.0)
	c.Check(nearestRank(input, 100), Equals, 100.0)
}

func (s *S) TestNearestRankFor101(c *C) {
	input := make([]float64, 101)

	for i := 0; i < 101; i++ {
		input[i] = float64(i + 1)
	}

	c.Check(nearestRank(input, 0), Equals, 1.0)
	c.Check(nearestRank(input, 50), Equals, 51.0)
	c.Check(nearestRank(input, 100), Equals, 101.0)
}

func (s *S) TestMedianReducer(c *C) {
	input := []float64{1, 2, 3}

	c.Check(MedianReducer(input), Equals, 2.0)
}

func (s *S) TestMinimum(c *C) {
	input := []float64{5, 1, 10, 1.1, 4}

	c.Check(MinimumReducer(input), Equals, 1.0)
}

func (s *S) TestMaximum(c *C) {
	input := []float64{5, 1, 10, 1.1, 4}

	c.Check(MaximumReducer(input), Equals, 10.0)
}
