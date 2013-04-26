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
	"time"
)

// An AccuracyAggregator encapsulates an action potential and
// records the average skew
type AccuracyAccumulator struct {
	ActionPotential
	AverageDelta time.Duration
	Count        int64
}

func NewAccuracyAccumulator(ap ActionPotential) *AccuracyAccumulator {
	return &AccuracyAccumulator{ap, 0, 0}
}

func (f *AccuracyAccumulator) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	now := time.Now()
	potential, fired := f.ActionPotential.AddPotentialAt(p, t)
	total_skew := int64(f.AverageDelta)*f.Count + int64(now.Sub(t))
	f.Count += 1
	f.AverageDelta = time.Duration(total_skew / f.Count)
	return potential, fired
}
