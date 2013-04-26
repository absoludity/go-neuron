// go-neuron - A neuron simulator for Go.
//
// Copyright (c) 2013 - Michael Nelson <absoludity@gmail.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package action_potential

import (
	"testing"
	"time"
)

func average(ds []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range ds {
		total += d
	}
	return total / time.Duration(len(ds))
}

func TestAccuracyAggretator(t *testing.T) {
	interval := 100 * time.Millisecond
	delays := []time.Duration{
		1 * interval,
		2 * interval,
		3 * interval,
	}
	accum := NewAccuracyAccumulator(new(Simple))

	for _, delay := range delays {
		accum.AddPotentialAt(5, time.Now().Add(-delay))
	}

	if accum.Count != int64(len(delays)) {
		t.Errorf("Expected Count=%d but was %d.", len(delays), accum.Count)
	}
	expected_average := average(delays)
	slop := 10 * time.Microsecond
	if accum.AverageDelta < expected_average-slop ||
		accum.AverageDelta > expected_average+slop {
		t.Errorf("Average delta was %s, expected [%s,%s]",
			accum.AverageDelta, expected_average-slop, expected_average+slop)
	}
}
