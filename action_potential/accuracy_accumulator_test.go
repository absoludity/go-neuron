package action_potential

import (
	"testing"
	"time"
)

func average(ds []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range ds {
		total += d
	}
	return total / time.Duration(len(ds))
}

func TestAccuracyAggretator(t *testing.T) {
	interval := 100 * time.Millisecond
	delays := []time.Duration{
		1 * interval,
		2 * interval,
		3 * interval,
	}
	accum := NewAccuracyAccumulator(new(Simple))

	for _, delay := range delays {
		accum.AddPotentialAt(5, time.Now().Add(-delay))
	}

	if accum.Count != int64(len(delays)) {
		t.Errorf("Expected Count=%d but was %d.", len(delays), accum.Count)
	}
	expected_average := average(delays)
	slop := 10 * time.Microsecond
	if accum.AverageDelta < expected_average-slop ||
		accum.AverageDelta > expected_average+slop {
		t.Errorf("Average delta was %s, expected [%s,%s]",
			accum.AverageDelta, expected_average-slop, expected_average+slop)
	}
}
