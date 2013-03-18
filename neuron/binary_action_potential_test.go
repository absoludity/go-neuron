package neuron

import (
	"testing"
	"time"
)

func verify(t *testing.T, testnum int, expected, result PotentialState) {
	if result != expected {
		t.Errorf("%d. Expected: %s, actual: %s", testnum, expected, result)
	}
}

var (
	now             = time.Now()
	before_decay    = now.Add(BINARY_DECAY_DURATION - time.Microsecond)
	after_decay     = now.Add(BINARY_DECAY_DURATION + time.Microsecond)
	before_inactive = now.Add(BINARY_ACTIVE_DURATION - time.Microsecond)
	after_inactive  = now.Add(BINARY_ACTIVE_DURATION + time.Microsecond)
	before_rest     = now.Add(BINARY_INACTIVE_DURATION - time.Microsecond)
	after_rest      = now.Add(BINARY_INACTIVE_DURATION + time.Microsecond)
)

var get_potential_cases = []struct {
	in  PotentialState
	at  time.Time
	out PotentialState
}{
	// A deactivated cell with a small potential will not have changed before the
	// decay time.
	{PotentialState{1, now, DEACTIVATED}, before_decay, PotentialState{1, now, DEACTIVATED}},
	// A deactivated cell with a potential will have returned to rest
	// potential after the decay duration.
	{PotentialState{1, now, DEACTIVATED}, after_decay, PotentialState{0, after_decay, DEACTIVATED}},
	// An activated cell will remain active for the active duration.
	{PotentialState{100, now, ACTIVATED}, before_inactive, PotentialState{PEAK_POTENTIAL, now, ACTIVATED}},
	// An activated cell will be inactive after the active duration.
	{PotentialState{100, now, ACTIVATED}, after_inactive, PotentialState{REFRACTORY_POTENTIAL, now.Add(BINARY_ACTIVE_DURATION), INACTIVATED}},
	// An inactivated cell will remain inactivated for the inactive duration.
	{PotentialState{REFRACTORY_POTENTIAL, now, INACTIVATED}, before_rest, PotentialState{REFRACTORY_POTENTIAL, now, INACTIVATED}},
	// An inactivated cell will switch back to deactivated after the
	// refractory period.
	{
		PotentialState{REFRACTORY_POTENTIAL, now, INACTIVATED},
		after_rest,
		PotentialState{REST_POTENTIAL, now.Add(BINARY_INACTIVE_DURATION),
			DEACTIVATED},
	},
}

func TestGetPotentialAt(t *testing.T) {
	for i, tt := range get_potential_cases {
		cb := BinaryActionPotential{tt.in}

		actual_potential := cb.GetPotentialAt(tt.at)

		verify(t, i, tt.out, cb.PotentialState)
		if actual_potential != tt.out.last_potential {
			t.Errorf("%d: Expected return value: %.1f, actual: %.1f",
				i, tt.out.last_potential, actual_potential)
		}
	}
}

var add_potential_cases = []struct {
	initial PotentialState
	at      time.Time
	in      Potential
	out     Potential
	fired   bool
	final   PotentialState
}{
	// A deactivated cell with no potential added is still deactivated.
	{
		PotentialState{0, now, DEACTIVATED},
		before_decay,
		0,
		0, false,
		PotentialState{0, now, DEACTIVATED},
	},
	// A deactivated  cell collects any added potential before decaying.
	{
		PotentialState{1, now, DEACTIVATED},
		before_decay,
		1,
		2, false,
		PotentialState{2, before_decay, DEACTIVATED},
	},
	// A deactivated cell's potential returns to rest without an active potential.
	{
		PotentialState{1, now, DEACTIVATED},
		after_decay,
		1,
		1, false,
		PotentialState{1, after_decay, DEACTIVATED},
	},
	// A deactivated cell that reaches threshold will reach the action potential.
	{
		PotentialState{1, now, DEACTIVATED},
		before_decay,
		THRESHOLD_POTENTIAL,
		PEAK_POTENTIAL, true,
		PotentialState{PEAK_POTENTIAL, before_decay, ACTIVATED},
	},
	// An activated cell ignores any further additions.
	{
		PotentialState{PEAK_POTENTIAL, now, ACTIVATED},
		before_inactive,
		1,
		PEAK_POTENTIAL, false,
		PotentialState{PEAK_POTENTIAL, now, ACTIVATED},
	},
	// A cell that switches to inactivated after an active
	// potential will be inhibited for a period also.
	{
		PotentialState{PEAK_POTENTIAL, now, ACTIVATED},
		after_inactive,
		1,
		REFRACTORY_POTENTIAL, false,
		PotentialState{REFRACTORY_POTENTIAL, now.Add(BINARY_ACTIVE_DURATION), INACTIVATED},
	},
	// An inactivated cell will, after the refractory period, accumulate
	// again.
	{
		PotentialState{REFRACTORY_POTENTIAL, now, INACTIVATED},
		after_rest,
		1,
		1, false,
		PotentialState{1, after_rest, DEACTIVATED},
	},
}

func TestAddPotentialAt(t *testing.T) {
	for i, tt := range add_potential_cases {
		cb := BinaryActionPotential{tt.initial}

		actual_potential, fired := cb.AddPotentialAt(tt.in, tt.at)

		if tt.fired != fired {
			t.Errorf("%d: Expected to fire: %s, actually fired: %s.",
				i, tt.fired, fired)
		}
		if tt.out != actual_potential {
			t.Errorf("%d: Expected return value: %.1f, actual: %.1f",
				i, tt.out, actual_potential)
		}
		verify(t, i, tt.final, cb.PotentialState)
	}
}
