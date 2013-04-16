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
