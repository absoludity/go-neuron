package neuron

import (
	"time"
)

// Any potential that doesn't reach the threshold and result
// in an action potential quickly decays. For simplicity the
// cell body starts with a binary decay function (ie. the
// potential disappears after Xms).
// TODO: Find a realistic value for this.
const (
	BINARY_DECAY_DURATION    = 3 * time.Millisecond
	BINARY_ACTIVE_DURATION   = 3 * time.Millisecond
	BINARY_INACTIVE_DURATION = 3 * time.Millisecond
)

type BinaryActionPotential struct {
	PotentialState
}

func (cb *BinaryActionPotential) GetPotentialAt(now time.Time) Potential {
	switch cb.state {
	case DEACTIVATED:
		decay_time := cb.last_change.Add(BINARY_DECAY_DURATION)
		if decay_time.Before(now) {
			cb.last_potential = 0
			cb.last_change = now
		}
	case ACTIVATED:
		inactive_time := cb.last_change.Add(BINARY_ACTIVE_DURATION)
		if inactive_time.Before(now) {
			cb.last_potential = REFRACTORY_POTENTIAL
			cb.state = INACTIVATED
			cb.last_change = inactive_time
		}
	case INACTIVATED:
		deactivated_time := cb.last_change.Add(BINARY_ACTIVE_DURATION)
		if deactivated_time.Before(now) {
			cb.state = DEACTIVATED
			cb.last_change = deactivated_time
			cb.last_potential = REST_POTENTIAL
		}
	}
	return cb.last_potential
}

func (cb *BinaryActionPotential) GetPotential() Potential {
	return cb.GetPotentialAt(time.Now())
}

func (cb *BinaryActionPotential) AddPotentialAt(potential Potential, now time.Time) (Potential, bool) {
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

func (cb *BinaryActionPotential) AddPotential(potential Potential) (Potential, bool) {
	return cb.AddPotentialAt(potential, time.Now())
}
