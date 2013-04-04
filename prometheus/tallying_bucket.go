// Copyright (c) 2013, Prometheus Team
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"fmt"
	"math"
	"sync"
)

const (
	lowerThird = 1.0 / 3.0
	upperThird = 2.0 * lowerThird
)

// A TallyingIndexEstimator is responsible for estimating the value of index for
// a given TallyingBucket, even though a TallyingBucket does not possess a
// collection of samples.  There are a few strategies listed below for how
// this value should be approximated.
type TallyingIndexEstimator func(minimum, maximum float64, index, observations int) float64

// Provide a filter for handling empty buckets.
func emptyFilter(e TallyingIndexEstimator) TallyingIndexEstimator {
	return func(minimum, maximum float64, index, observations int) float64 {
		if observations == 0 {
			return math.NaN()
		}

		return e(minimum, maximum, index, observations)
	}
}

var (
	minimumEstimator = emptyFilter(func(minimum, maximum float64, _, observations int) float64 {
		return minimum
	})

	maximumEstimator = emptyFilter(func(minimum, maximum float64, _, observations int) float64 {
		return maximum
	})

	averageEstimator = emptyFilter(func(minimum, maximum float64, _, observations int) float64 {
		return AverageReducer([]float64{minimum, maximum})
	})

	uniformEstimator = emptyFilter(func(minimum, maximum float64, index, observations int) float64 {
		if observations == 1 {
			return minimum
		}

		location := float64(index) / float64(observations)

		if location > upperThird {
			return maximum
		} else if location < lowerThird {
			return minimum
		}

		return AverageReducer([]float64{minimum, maximum})
	})
)

// These are the canned TallyingIndexEstimators.
var (
	// Report the smallest observed value in the bucket.
	MinimumEstimator = minimumEstimator
	// Report the largest observed value in the bucket.
	MaximumEstimator = maximumEstimator
	// Report the average of the extrema.
	AverageEstimator = averageEstimator
	// Report the minimum value of the index if it is in the lower-third of
	// observations, the average if in the middle-third, and the maximum if in
	// the largest third
	UniformEstimator = uniformEstimator
)

// A TallyingBucket is a Bucket that tallies when an object is added to it.
// Upon insertion, an object is compared against collected extrema and noted
// as a new minimum or maximum if appropriate.
type TallyingBucket struct {
	estimator        TallyingIndexEstimator
	largestObserved  float64
	mutex            sync.RWMutex
	observations     int
	smallestObserved float64
}

func (b *TallyingBucket) Add(value float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.observations += 1
	b.smallestObserved = math.Min(value, b.smallestObserved)
	b.largestObserved = math.Max(value, b.largestObserved)
}

func (b TallyingBucket) String() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	observations := b.observations

	if observations == 0 {
		return fmt.Sprintf("[TallyingBucket (Empty)]")
	}

	return fmt.Sprintf("[TallyingBucket (%f, %f); %d items]", b.smallestObserved, b.largestObserved, observations)
}

func (b TallyingBucket) Observations() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.observations
}

func (b TallyingBucket) ValueForIndex(index int) float64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.estimator(b.smallestObserved, b.largestObserved, index, b.observations)
}

func (b *TallyingBucket) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.largestObserved = math.SmallestNonzeroFloat64
	b.observations = 0
	b.smallestObserved = math.MaxFloat64
}

// Produce a TallyingBucket with sane defaults.
func DefaultTallyingBucket() TallyingBucket {
	return TallyingBucket{
		estimator:        MinimumEstimator,
		largestObserved:  math.SmallestNonzeroFloat64,
		smallestObserved: math.MaxFloat64,
	}
}

func CustomTallyingBucket(estimator TallyingIndexEstimator) TallyingBucket {
	return TallyingBucket{
		estimator:        estimator,
		largestObserved:  math.SmallestNonzeroFloat64,
		smallestObserved: math.MaxFloat64,
	}
}

// This is used strictly for testing.
func tallyingBucketBuilder() Bucket {
	b := DefaultTallyingBucket()
	return &b
}
