package neuron

import (
	"time"
)

type ActivationEvent struct {
	Time   time.Time
	Neuron *Neuron
}

type ActivationStream chan ActivationEvent

type Axon struct {
	Terminals []Neuron
	Delay     time.Duration
}

type Neuron struct {
	Axon             Axon
	ActivationStream ActivationStream
	ActionPotential
}

func (n *Neuron) AddPotentialAt(p Potential, t time.Time) Potential {
	potential, fired := n.ActionPotential.AddPotentialAt(p, t)
	if fired {
		n.ActivationStream <- ActivationEvent{t, n}
	}
	return potential
}
