package neuron

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
	fake := action_potential.NewEventRecorder(new(action_potential.Simple))
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

	if len(fake.Events) != len(delays) {
		t.Errorf("Expected %d calls to AddPotential, received %d.",
			len(delays), len(fake.Events))
	}
	for i, delay := range delays {
		expected := action_potential.AddPotentialEvent{5, now.Add(delay)}
		if fake.Events[i] != expected {
			t.Errorf("Expected %s, got %s.", expected, fake.Events[i])
		}
	}
}

func TestOrdersAccordingToDelay(t *testing.T) {
	as := make(ActivationStream, 2)
	now := time.Now()
	fake := action_potential.NewEventRecorder(new(action_potential.Simple))
	delays := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
	}
	as <- ActivationEvent{now, makeNeuronWithTerminal(fake, 5*time.Millisecond)}
	as <- ActivationEvent{now, makeNeuronWithTerminal(fake, 1*time.Millisecond)}
	close(as)

	as.Process()

	if len(fake.Events) != 2 {
		t.Errorf("Expected 2 calls to AddPotential, received %d.",
			len(delays), len(fake.Events))
	}
	expected := action_potential.AddPotentialEvent{5, now.Add(1 * time.Millisecond)}
	if fake.Events[0] != expected {
		t.Errorf("Expected first call to be %s, but was %s",
			expected, fake.Events[0])
	}
	expected = action_potential.AddPotentialEvent{5, now.Add(5 * time.Millisecond)}
	if fake.Events[1] != expected {
		t.Errorf("Expected second call to be %s, but was %s",
			expected, fake.Events[1])
	}
}
