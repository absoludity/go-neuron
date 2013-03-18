package neuron

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
