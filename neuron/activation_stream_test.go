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

func TestProcessUntilEmpty(t *testing.T) {
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

	as.ProcessUntilEmpty()

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

func TestActivationStreamAccuracy(t *testing.T) {
	accum := action_potential.NewAccuracyAccumulator(
		new(action_potential.Simple))
	activation_stream := make(ActivationStream, 1000)
	neurons := make([]*Neuron, 1000)
	// Each neuron has a delay ranging from 10ms to 1009ms
	for i := 0; i < 1000; i++ {
		neurons[i] = &Neuron{
			Axon{
				[]action_potential.ActionPotential{accum},
				time.Duration(i+10) * time.Millisecond,
			},
			&activation_stream,
			action_potential.NewAlwaysFirer(new(action_potential.Simple)),
		}
	}
	for _, n := range neurons {
		n.AddPotential(0)
	}

	activation_stream.ProcessUntilEmpty()

	expected_accuracy := time.Duration(10) * time.Microsecond
	if accum.AverageDelta*accum.AverageDelta > expected_accuracy*expected_accuracy {
		t.Errorf("Expected an accuracy better than %s, but "+
			"average delta was %s.", expected_accuracy, accum.AverageDelta)
	}
}

func TestActivationStreamDelay(t *testing.T) {
	activation_stream := make(ActivationStream, 1000)
	event_recorder := action_potential.NewEventRecorder(
		new(action_potential.Simple))
	axon_delay := time.Duration(1) * time.Microsecond
	end := &Neuron{
		Axon{
			[]action_potential.ActionPotential{event_recorder},
			axon_delay,
		},
		&activation_stream,
		action_potential.NewAlwaysFirer(new(action_potential.Simple)),
	}
	prev := end

	// A string of neurons together with 1us axon delays.
	for i := 0; i < 999; i++ {
		neuron := &Neuron{
			Axon{
				[]action_potential.ActionPotential{prev},
				axon_delay,
			},
			&activation_stream,
			action_potential.NewAlwaysFirer(new(action_potential.Simple)),
		}
		prev = neuron
	}
	start := prev

	started_at := time.Now()
	start.AddPotentialAt(0, started_at)
	go activation_stream.Process()

	time.Sleep(axon_delay * 1000)
	close(activation_stream)
	if len(event_recorder.Events) != 1 {
		t.Errorf("Expected 1 event, got %d.", len(event_recorder.Events))
	}

	delay := event_recorder.Events[0].Time.Sub(started_at)
	t.Errorf("Delay was %s.", delay)
	// ^^ Always exactly 1ms as it's the time it should be added, not the time
	// it was actually added.
}
