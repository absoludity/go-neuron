package neuron

import (
	"fmt"
	"testing"
	"time"
)

type ps struct {
	potential float32
	state     NeuronState
}

func (p ps) String() string {
	return fmt.Sprintf("(%.1f, %s)", p.potential, p.state)
}

func verify(t *testing.T, testnum int, expected, output ps) {
	if output != expected {
		t.Errorf("%d. Expected: %s, actual: %s", testnum, expected, output)
	}
}

func TestGetPotentialAt(t *testing.T) {
	cases := []struct {
		in       ps
		duration time.Duration
		out      ps
	}{
		{ps{1, DEACTIVATED}, BINARY_DECAY_DURATION - time.Microsecond, ps{1, DEACTIVATED}},
		{ps{1, DEACTIVATED}, BINARY_DECAY_DURATION + time.Microsecond, ps{0, DEACTIVATED}},
	}

	now := time.Now()
	for i, tt := range cases {
		cb := CellBody{tt.in.potential, now, tt.in.state}

		actual_potential := cb.GetPotentialAt(now.Add(tt.duration))

		verify(t, i, tt.out, ps{actual_potential, cb.state})
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
		cb := CellBody{tt.in, time.Now().Add(-tt.duration), DEACTIVATED}

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
		cb := CellBody{tt.last_potential, now.Add(-tt.since_last_change), DEACTIVATED}

		actual_potential := cb.AddPotentialAt(tt.in, now)

		if actual_potential != tt.out {
			t.Errorf("%d. Expected potential:%f, actual: %f.", i, tt.out,
				actual_potential)
		}
	}
}
