// Copyright (c) 2013, Prometheus Team
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"bytes"
	"container/heap"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type AccumulatingBucket struct {
	elements       priorityQueue
	evictionPolicy EvictionPolicy
	maximumSize    int
	mutex          sync.RWMutex
	observations   int
}

// AccumulatingBucketBuilder is a convenience method for generating a
// BucketBuilder that produces AccumatingBucket entries with a certain
// behavior set.
func AccumulatingBucketBuilder(evictionPolicy EvictionPolicy, maximumSize int) BucketBuilder {
	return func() Bucket {
		return &AccumulatingBucket{
			elements:       make(priorityQueue, 0, maximumSize),
			evictionPolicy: evictionPolicy,
			maximumSize:    maximumSize,
		}
	}
}

// Add a value to the bucket.  Depending on whether the bucket is full, it may
// trigger an eviction of older items.
func (b *AccumulatingBucket) Add(value float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.observations++
	size := len(b.elements)

	v := item{
		Priority: -1 * time.Now().UnixNano(),
		Value:    value,
	}

	if size == b.maximumSize {
		b.evictionPolicy(&b.elements)
	}

	heap.Push(&b.elements, &v)
}

func (b AccumulatingBucket) String() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	buffer := &bytes.Buffer{}

	fmt.Fprintf(buffer, "[AccumulatingBucket with %d elements and %d capacity] { ", len(b.elements), b.maximumSize)

	for i := 0; i < len(b.elements); i++ {
		fmt.Fprintf(buffer, "%f, ", b.elements[i].Value)
	}

	fmt.Fprintf(buffer, "}")

	return buffer.String()
}

func (b AccumulatingBucket) ValueForIndex(index int) float64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	elementCount := len(b.elements)

	if elementCount == 0 {
		return math.NaN()
	}

	sortedElements := make([]float64, elementCount)

	for i, element := range b.elements {
		sortedElements[i] = element.Value.(float64)
	}

	sort.Float64s(sortedElements)

	// N.B.(mtp): Interfacing components should not need to comprehend what
	//            eviction and storage container strategies used; therefore,
	//            we adjust this silently.
	targetIndex := int(float64(elementCount-1) * (float64(index) / float64(b.observations)))

	return sortedElements[targetIndex]
}

func (b AccumulatingBucket) Observations() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.observations
}

func (b *AccumulatingBucket) Reset() {
	b.mutex.Lock()
	defer b.mutex.RUnlock()

	for i := 0; i < b.elements.Len(); i++ {
		b.elements.Pop()
	}

	b.observations = 0
}
