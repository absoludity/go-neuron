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
/*
	Package action_potential provides an interface and simple
	implementation of the short-lasting electrical event occuring
	in neurons.

	http://en.wikipedia.org/wiki/Action_potential
*/
package action_potential

import (
	"fmt"
	"time"
)

type Potential float32

type ActionPotential interface {
	GetPotential() Potential
	GetPotentialAt(time.Time) Potential
	AddPotential(Potential) (Potential, bool)
	AddPotentialAt(Potential, time.Time) (Potential, bool)
}

// Typically 15mV above the resting potential.
// Move these into a subclass perhaps, so it's possible
// to have different classes of neurons without requiring
// each neuron to include the values.
const (
	REST_POTENTIAL       Potential = 0
	THRESHOLD_POTENTIAL  Potential = 15
	PEAK_POTENTIAL       Potential = 100
	REFRACTORY_POTENTIAL Potential = -15
)

type ActivationState int

// An ActionPotential can be deactivated, activated or
// inactivated.
const (
	DEACTIVATED ActivationState = iota
	ACTIVATED
	INACTIVATED
)

func (as ActivationState) String() string {
	switch as {
	case DEACTIVATED:
		return "Deactivated"
	case ACTIVATED:
		return "Activated"
	case INACTIVATED:
		return "Inactivated"
	}
	return "Unknown"
}

// PotentialState stores the data required to determine
// a potential at a given time (internally the state,
// the previous potential and the time at which the potential
// last changed.)
type PotentialState struct {
	// The zero value is used as the resting potential.
	last_potential Potential
	last_change    time.Time
	state          ActivationState
}

func (ps PotentialState) String() string {
	return fmt.Sprintf("%s (%.1f since %s ago)",
		ps.state, ps.last_potential, time.Now().Sub(ps.last_change))
}
