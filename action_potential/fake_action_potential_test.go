package action_potential

import (
	"testing"
	"time"
)

func TestFakeAddPotentialAt(t *testing.T) {
	now = time.Now()
	events := []AddPotentialEvent{
		{1.1, now},
		{2.2, now.Add(time.Hour * 24)},
		{3.3, now.Add(time.Hour * 48)},
	}
	fake := NewFakeActionPotential()

	for _, e := range events {
		fake.AddPotentialAt(e.Potential, e.Time)
	}

	if len(fake) != 3 {
		t.Errorf("Expected %d events, received %d.", len(events), len(fake))
	}
	for i, e := range fake {
		if e != events[i] {
			t.Errorf("Expected fake event %s to be %s, but was %s.",
				i, events[i], e)

		}
	}
}
