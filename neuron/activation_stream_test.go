package neuron

// Add 3 events + termination to activation stream then call Process()
// Verify the additions at the correct times.

import (
	"github.com/absoludity/go-neuron/action_potential"
	"testing"
	"time"
)

// makeNeuronWithTerminal returns a pointer to a neuron with the given
// neuron and delay added to the axon terminal.
func makeNeuronWithTerminal(ap action_potential.ActionPotential, delay time.Duration) *Neuron {
	return &Neuron{
		Axon{
			[]action_potential.ActionPotential{ap},
			delay,
		},
		nil,
		nil,
	}
}

func TestProcess5Simultaneous(t *testing.T) {
	as := make(ActivationStream, 5)
	now := time.Now()
	fake := action_potential.NewFakeActionPotential()
	delays := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
	}
	for _, delay := range delays {
		as <- ActivationEvent{now, makeNeuronWithTerminal(fake, delay)}
	}
	close(as)

	as.Process()

	if len(fake) != len(delays) {
		t.Errorf("Expected %d calls to AddPotential, received %d.",
			len(delays), len(fake))
	}
	for i, delay := range delays {
		expected := action_potential.AddPotentialEvent{5, now.Add(delay)}
		if fake[i] != expected {
			t.Errorf("Expected %s, got %s.", expected, fake[i])
		}
	}
}

// func TestOrdersAccordingToDelay(t *testing.T) {
// 	as := make(ActivationStream, 2)
// 	now := time.Now()
// 	fake := action_potential.NewFakeActionPotential()
// 	delays := []time.Duration{
// 		1 * time.Millisecond,
// 		2 * time.Millisecond,
// 		3 * time.Millisecond,
// 		4 * time.Millisecond,
// 		5 * time.Millisecond,
// 	}
// 	as <- ActivationEvent{
// 		now,
// 		&Neuron{
// 			Axon{
// 				[]action_potential.ActionPotential{&fake},
// 				1 * time.Millisecond,
// 			},
// 			nil,
// 			nil,
// 		},
// 	}
// 	close(as)

// 	as.Process()

// 	if len(fake) != len(delays) {
// 		t.Errorf("Expected %d calls to AddPotential, received %d.",
// 			len(delays), len(fake))
// 	}
// 	for i, delay := range delays {
// 		expected := action_potential.AddPotentialEvent{5, now.Add(delay)}
// 		if fake[i] != expected {
// 			t.Errorf("Expected %s, got %s.", expected, fake[i])
// 		}
// 	}

// }
