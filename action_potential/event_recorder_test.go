package action_potential

import (
	"testing"
	"time"
)

func TestEventRecorderAddPotentialAt(t *testing.T) {
	now = time.Now()
	events := []AddPotentialEvent{
		{1.0, now},
		{2.0, now.Add(time.Microsecond * 1)},
		{3.0, now.Add(time.Microsecond * 2)},
	}
	fake := NewEventRecorder(new(Simple))

	var (
		actual_potential Potential
		fired            bool
	)
	for _, e := range events {
		actual_potential, fired = fake.AddPotentialAt(e.Potential, e.Time)
	}

	if len(fake.Events) != 3 {
		t.Errorf("Expected %d events, received %d.", len(events), len(fake.Events))
	}
	for i, e := range fake.Events {
		if e != events[i] {
			t.Errorf("Expected fake event %s to be %s, but was %s.",
				i, events[i], e)
		}
	}
	if fired {
		t.Error("Unexpected firing of action potential.")
	}
	if actual_potential != 6.0 {
		t.Errorf("Expected return value: 6.0, actual %.1f.", actual_potential)
	}
}
