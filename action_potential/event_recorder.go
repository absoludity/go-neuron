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

type AddPotentialEvent struct {
	Potential Potential
	Time      time.Time
	RealTime  time.Time
}

// An EventRecorder encapsulates an action potential and records
// each AddPotential event that occurs.
type EventRecorder struct {
	ActionPotential
	Events []AddPotentialEvent
}

func NewEventRecorder(ap ActionPotential) *EventRecorder {
	return &EventRecorder{ap, make([]AddPotentialEvent, 0, 10)}
}

func (f *EventRecorder) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	f.Events = append(f.Events, AddPotentialEvent{p, t, time.Now()})
	return f.ActionPotential.AddPotentialAt(p, t)
}

func (f *EventRecorder) AddPotential(p Potential) (Potential, bool) {
	return f.AddPotentialAt(p, time.Now())
}
