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

// The Simple activates for a duration
// when the initial threshold is reached, then
// is inactive for a duration before switching back to deactivated.
const (
	SIMPLE_DECAY_DURATION    = 3 * time.Millisecond
	SIMPLE_ACTIVE_DURATION   = 3 * time.Millisecond
	SIMPLE_INACTIVE_DURATION = 3 * time.Millisecond
)

// The Simple is a simple implementation of
// the action potential interface.
type Simple struct {
	PotentialState
}

// GetPotentialAt determines and returns the potential at a given
// point in time.
func (cb *Simple) GetPotentialAt(now time.Time) Potential {
	switch cb.state {
	case DEACTIVATED:
		decay_time := cb.last_change.Add(SIMPLE_DECAY_DURATION)
		if decay_time.Before(now) {
			cb.last_potential = 0
			cb.last_change = now
		}
	case ACTIVATED:
		inactive_time := cb.last_change.Add(SIMPLE_ACTIVE_DURATION)
		if inactive_time.Before(now) {
			cb.last_potential = REFRACTORY_POTENTIAL
			cb.state = INACTIVATED
			cb.last_change = inactive_time
		}
	case INACTIVATED:
		deactivated_time := cb.last_change.Add(SIMPLE_ACTIVE_DURATION)
		if deactivated_time.Before(now) {
			cb.state = DEACTIVATED
			cb.last_change = deactivated_time
			cb.last_potential = REST_POTENTIAL
		}
	}
	return cb.last_potential
}

// GetPotential determines and returns the potential at the time it
// is called.
func (cb *Simple) GetPotential() Potential {
	return cb.GetPotentialAt(time.Now())
}

// AddPotentialAt adds the specified potential based on the existing
// potential at the specified time.
func (cb *Simple) AddPotentialAt(potential Potential, now time.Time) (Potential, bool) {
	prev := cb.last_potential
	fired := false
	current_potential := cb.GetPotentialAt(now)
	switch cb.state {
	case DEACTIVATED:
		cb.last_potential = current_potential + potential
		if cb.last_potential > THRESHOLD_POTENTIAL {
			cb.state = ACTIVATED
			cb.last_potential = PEAK_POTENTIAL
			fired = true
		}
		if cb.last_potential != prev {
			cb.last_change = now
		}
	}
	return cb.last_potential, fired
}

// AddPotential adds the specified potential based on the existing
// potential at the time it is called.
func (cb *Simple) AddPotential(potential Potential) (Potential, bool) {
	return cb.AddPotentialAt(potential, time.Now())
}
