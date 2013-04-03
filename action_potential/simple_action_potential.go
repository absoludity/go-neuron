package action_potential

import (
	"time"
)

// The SimpleActionPotential activates for a duration
// when the initial threshold is reached, then
// is inactive for a duration before switching back to deactivated.
const (
	SIMPLE_DECAY_DURATION    = 3 * time.Millisecond
	SIMPLE_ACTIVE_DURATION   = 3 * time.Millisecond
	SIMPLE_INACTIVE_DURATION = 3 * time.Millisecond
)

// The SimpleActionPotential is a simple implementation of
// the action potential interface.
type SimpleActionPotential struct {
	PotentialState
}

// Creates and returns a pointer to an action potential
// with the given potential at the specified time.
// Rename (New too general in module scope) and move to unexported
// test function, if only used there.
func New(p Potential, t time.Time) ActionPotential {
	return &SimpleActionPotential{PotentialState{p, t, DEACTIVATED}}
}

// GetPotentialAt determines and returns the potential at a given
// point in time.
func (cb *SimpleActionPotential) GetPotentialAt(now time.Time) Potential {
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
func (cb *SimpleActionPotential) GetPotential() Potential {
	return cb.GetPotentialAt(time.Now())
}

// AddPotentialAt adds the specified potential based on the existing
// potential at the specified time.
func (cb *SimpleActionPotential) AddPotentialAt(potential Potential, now time.Time) (Potential, bool) {
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
func (cb *SimpleActionPotential) AddPotential(potential Potential) (Potential, bool) {
	return cb.AddPotentialAt(potential, time.Now())
}
