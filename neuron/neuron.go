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

// An ActivationEvent records the neuron and time at which it
// was activated.
type ActivationEvent struct {
	Time   time.Time
	Neuron *Neuron
}

// An ActivationStream communicates the activation events for further
// processing.
type ActivationStream chan ActivationEvent

func (as *ActivationStream) Process() {
	for {
		ae := <-*as
		if ae.Neuron == nil {
			return
		}
		axon := ae.Neuron.Axon
		for _, n := range axon.Terminals {
			// Should the potential for each be relative to total
			// potential, or constant, or divided by the num of terminals?
			n.AddPotentialAt(5.0, ae.Time.Add(axon.Delay))
		}

	}
}

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
func (n *Neuron) AddPotentialAt(p action_potential.Potential, t time.Time) action_potential.Potential {
	potential, fired := n.ActionPotential.AddPotentialAt(p, t)
	if fired {
		*n.ActivationStream <- ActivationEvent{t, n}
	}
	return potential
}
