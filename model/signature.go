// Copyright (c) 2013, Prometheus Team
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"hash/fnv"
	"sort"
)

// cache the signature of an empty label set.
var emptyLabelSignature = fnv.New64a().Sum64()

// LabelsToSignature provides a way of building a unique signature
// (i.e., fingerprint) for a given label set sequence.
func LabelsToSignature(labels map[string]string) uint64 {
	if len(labels) == 0 {
		return emptyLabelSignature
	}

	names := make([]string, 0, len(labels))
	for name := range labels {
		names = append(names, name)
	}

	sort.Strings(names)

	hasher := fnv.New64a()

	for _, name := range names {
		hasher.Write([]byte(name))
		hasher.Write([]byte(labels[name]))
	}

	return hasher.Sum64()
}

// LabelValuesToSignature provides a way of building a unique signature
// (i.e., fingerprint) for a given set of label's values.
func LabelValuesToSignature(labels map[string]string) uint64 {
	if len(labels) == 0 {
		return emptyLabelSignature
	}

	names := make([]string, 0, len(labels))
	for name := range labels {
		names = append(names, name)
	}

	sort.Strings(names)

	hasher := fnv.New64a()

	for _, name := range names {
		hasher.Write([]byte(labels[name]))
	}

	return hasher.Sum64()
}
