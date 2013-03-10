package neuron

import (
	"testing"
	"time"
)

func TestGetPotentialAt(t *testing.T) {
	cases := []struct {
		in       float32
		duration time.Duration
		out      float32
	}{
		{1, BINARY_DECAY_DURATION - time.Microsecond, 1},
		{1, (BINARY_DECAY_DURATION + time.Microsecond), 0},
	}

	now := time.Now()
	for i, tt := range cases {
		cb := CellBody{tt.in, now}

		actual_potential := cb.GetPotentialAt(now.Add(tt.duration))

		if actual_potential != tt.out {
			t.Errorf("%d. Expected potential: %f, actual: %f.",
				i, tt.out, actual_potential)
		}
	}
}

func TestGetPotential(t *testing.T) {
	cases := []struct {
		in       float32
		duration time.Duration
		out      float32
	}{
		{1, 0, 1},
		{1, (BINARY_DECAY_DURATION + time.Microsecond), 0},
	}

	for i, tt := range cases {
		cb := CellBody{tt.in, time.Now().Add(-tt.duration)}

		actual_potential := cb.GetPotential()

		if actual_potential != tt.out {
			t.Errorf("%d. Expected potential: %f, actual: %f.",
				i, tt.out, actual_potential)
		}
	}
}

func TestAddPotentialAt(t *testing.T) {
	cases := []struct {
		last_potential    float32
		since_last_change time.Duration
		in                float32
		out               float32
	}{
		{0, BINARY_DECAY_DURATION - time.Millisecond, 0, 0},
		{1, BINARY_DECAY_DURATION - time.Millisecond, 1, 2},
		{1, BINARY_DECAY_DURATION + time.Microsecond, 1, 1},
	}
	for i, tt := range cases {
		now := time.Now()
		cb := CellBody{tt.last_potential, now.Add(-tt.since_last_change)}

		actual_potential := cb.AddPotentialAt(tt.in, now)

		if actual_potential != tt.out {
			t.Errorf("%d. Expected potential:%f, actual: %f.", i, tt.out,
				actual_potential)
		}
	}
}
