// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package promauto

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestWrapNil(t *testing.T) {
	// A nil registerer should be treated as a no-op by promauto, even when wrapped.
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"foo": "bar"}, nil)
	c := With(registerer).NewCounter(prometheus.CounterOpts{Name: "test"})
	c.Inc()
}
