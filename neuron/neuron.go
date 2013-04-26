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
	Package neuron implements a time-based neuron simulator.

	A neuron implements a simple interface that responds to additional
	potential being added at a point in time. If the additional potential
	over time reaches a threshold, the neuron fires and sends an activation
	event.
*/
package neuron

import (
	"github.com/absoludity/go-neuron/action_potential"
	"time"
)

// An Axon can have many terminals connecting to other neurons and
// an associated delay between the neurons activation and when the
// signal reaches the terminals.
type Axon struct {
	Terminals []action_potential.ActionPotential
	Delay     time.Duration
}

// A neuron itself is an ActionPotential implementation,
// together with a single Axon and an activation stream with which
// signals are communicated.
type Neuron struct {
	Axon             Axon
	ActivationStream *ActivationStream
	action_potential.ActionPotential
}

// AddPotentialAt updates the default implementation provided by
// the embedded ActivationPotential ensuring that any resulting activation
// is communicated to the stream.
func (n *Neuron) AddPotentialAt(p action_potential.Potential, t time.Time) (action_potential.Potential, bool) {
	potential, fired := n.ActionPotential.AddPotentialAt(p, t)
	if fired {
		*n.ActivationStream <- ActivationEvent{t, n}
	}
	return potential, fired
}

func (n *Neuron) AddPotential(p action_potential.Potential) (action_potential.Potential, bool) {
	return n.AddPotentialAt(p, time.Now())
}
