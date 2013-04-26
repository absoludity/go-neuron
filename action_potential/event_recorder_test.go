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

type EventRecorderCase struct {
	Potential Potential
	Time      time.Time
}

func TestEventRecorderAddPotentialAt(t *testing.T) {
	now = time.Now()
	events := []EventRecorderCase{
		{1.0, now},
		{2.0, now.Add(time.Microsecond * 1)},
		{3.0, now.Add(time.Microsecond * 2)},
	}
	fake := NewEventRecorder(new(Simple))

	var (
		actual_potential Potential
		fired            bool
	)
	for _, e := range events {
		actual_potential, fired = fake.AddPotentialAt(e.Potential, e.Time)
	}

	if len(fake.Events) != 3 {
		t.Errorf("Expected %d events, received %d.", len(events), len(fake.Events))
	}
	for i, e := range fake.Events {
		if e.Potential != events[i].Potential ||
			e.Time != events[i].Time {
			t.Errorf("Expected fake event %s to be %s, but was %s.",
				i, events[i], e)
		}
	}
	if fired {
		t.Error("Unexpected firing of action potential.")
	}
	if actual_potential != 6.0 {
		t.Errorf("Expected return value: 6.0, actual %.1f.", actual_potential)
	}
}

func TestEventRecorderAddPotential(t *testing.T) {
	potentials := []Potential{1.0, 2.0, 3.0}
	fake := NewEventRecorder(new(Simple))

	var (
		actual_potential Potential
		fired            bool
	)
	start := time.Now()
	for _, p := range potentials {
		actual_potential, fired = fake.AddPotential(p)
	}
	end := time.Now()

	if len(fake.Events) != 3 {
		t.Errorf("Expected %d events, received %d.", len(potentials), len(fake.Events))
	}
	previous_time := start
	for i, e := range fake.Events {
		if e.Potential != potentials[i] {
			t.Errorf("Expected fake event %s potential to be %s, but was %s.",
				i, potentials[i], e.Potential)
		}
		if e.Time.Before(previous_time) {
			t.Errorf("Expected fake events to be ordered "+
				"(event time %s <= previous time: %s", e.Time, previous_time)
		}
		previous_time = e.Time
	}
	if fake.Events[len(fake.Events)-1].Time.After(end) {
		t.Errorf("Last event's timestamp is after processing ended.")
	}
	if fired {
		t.Error("Unexpected firing of action potential.")
	}
	if actual_potential != 6.0 {
		t.Errorf("Expected return value: 6.0, actual %.1f.", actual_potential)
	}
}
