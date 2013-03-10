package neuron

import (
	"time"
)

type CellBody struct {
	// The zero value is used as the resting potential.
	last_potential float32
	last_change    time.Time
}

// Typically 15mV above the resting potential.
// Move these into a subclass perhaps, so it's possible
// to have different classes of neurons without requiring
// each neuron to include the values.
const THRESHHOLD_POTENTIAL = 15

// Any potential that doesn't reach the threshold and result
// in an action potential quickly decays. For simplicity the
// cell body starts with a binary decay function (ie. the
// potential disappears after Xms).
// TODO: Find a realistic value for this.
const BINARY_DECAY_DURATION = 3 * time.Millisecond

func (cb *CellBody) GetPotentialAt(now time.Time) float32 {
	if now.Sub(cb.last_change) > BINARY_DECAY_DURATION {
		cb.last_potential = 0
		cb.last_change = now
	}
	return cb.last_potential
}

func (cb *CellBody) GetPotential() float32 {
	return cb.GetPotentialAt(time.Now())
}

func (cb *CellBody) AddPotentialAt(potential float32, now time.Time) float32 {
	cb.last_potential = cb.GetPotentialAt(now) + potential
	cb.last_change = now
	return cb.last_potential
}

func (cb *CellBody) AddPotential(potential float32) float32 {
	return cb.AddPotentialAt(potential, time.Now())
}

// type FireEvent struct{
//     Connections Connections
//     float Charge
// }

// type Connections map[*CellBody]float
