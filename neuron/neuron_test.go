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
package neuron

import (
	"github.com/absoludity/go-neuron/action_potential"
	"testing"
	"time"
)

func TestNeuronFire(t *testing.T) {
	now := time.Now()
	cb := new(action_potential.Simple)
	cb.AddPotentialAt(1, now)
	as := make(ActivationStream, 1)
	n := &Neuron{Axon{}, &as, cb}
	at := now.Add(time.Microsecond * 5)

	actual_potential, fired := n.AddPotentialAt(action_potential.THRESHOLD_POTENTIAL, at)

	if actual_potential != action_potential.PEAK_POTENTIAL {
		t.Error("Expected potential:", action_potential.PEAK_POTENTIAL,
			"Actual potential: ", actual_potential)
	}
	if !fired {
		t.Errorf("Expected action potential to fire.")
	}
	ae, ok := <-as
	if !ok {
		t.Errorf("No activation event received on activation stream.")
	}
	if n != ae.Neuron {
		t.Errorf("Expected activation event for %s, but received %s",
			n, ae.Neuron)
	}
	if at != ae.Time {
		t.Errorf("Received activation event for incorrect time.")
	}
}

func TestNeuronNoFire(t *testing.T) {
	now := time.Now()
	ap := new(action_potential.Simple)
	ap.AddPotentialAt(1, now)
	as := make(ActivationStream, 1)
	n := &Neuron{Axon{}, &as, ap}
	at := now.Add(time.Microsecond * 5)

	actual_potential, fired := n.AddPotentialAt(5, at)

	if actual_potential != 6 {
		t.Error("Expected potential:", 6,
			"Actual potential: ", actual_potential)
	}
	if fired {
		t.Errorf("Expected action potential not to fire.")
	}

	// Close the activation stream so it doesn't block when empty.
	close(as)
	ae := <-as
	nil_activation_event := ActivationEvent{}
	if ae != nil_activation_event {
		t.Error("Expected the activation stream to be empty, but found ",
			ae)
	}
}

func TestNeuronFireAddPotential(t *testing.T) {
	now := time.Now()
	cb := new(action_potential.Simple)
	cb.AddPotentialAt(1, now)
	as := make(ActivationStream, 1)
	n := &Neuron{Axon{}, &as, cb}

	actual_potential, fired := n.AddPotential(action_potential.THRESHOLD_POTENTIAL)

	if actual_potential != action_potential.PEAK_POTENTIAL {
		t.Error("Expected potential:", action_potential.PEAK_POTENTIAL,
			"Actual potential: ", actual_potential)
	}
	if !fired {
		t.Errorf("Expected action potential to fire.")
	}
	ae, ok := <-as
	if !ok {
		t.Errorf("No activation event received on activation stream.")
	}
	if n != ae.Neuron {
		t.Errorf("Expected activation event for %s, but received %s",
			n, ae.Neuron)
	}
}
