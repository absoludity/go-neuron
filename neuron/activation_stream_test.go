package neuron

import (
	"github.com/absoludity/go-neuron/action_potential"
	"math"
	"testing"
	"time"
)

// makeNeuronWithTerminal returns a pointer to a neuron with the given
// neuron and delay added to the axon terminal.
func makeNeuronWithTerminal(terminal action_potential.ActionPotential,
	delay time.Duration, as *ActivationStream, ap action_potential.ActionPotential) *Neuron {
	return &Neuron{
		Axon{
			[]action_potential.ActionPotential{terminal},
			delay,
		},
		as,
		ap,
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
		as <- ActivationEvent{now, makeNeuronWithTerminal(fake, delay, nil, nil)}
	}
	close(as)

	as.Process()

	if len(fake.Events) != len(delays) {
		t.Errorf("Expected %d calls to AddPotential, received %d.",
			len(delays), len(fake.Events))
	}
	for i, delay := range delays {
		expected_time := now.Add(delay)
		if fake.Events[i].Time != expected_time {
			t.Errorf("Expected potential adaded at %s, got %s.",
				expected_time, fake.Events[i].Time)
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
	as <- ActivationEvent{now, makeNeuronWithTerminal(fake, 5*time.Millisecond, nil, nil)}
	as <- ActivationEvent{now, makeNeuronWithTerminal(fake, 1*time.Millisecond, nil, nil)}
	close(as)

	as.Process()

	if len(fake.Events) != 2 {
		t.Errorf("Expected 2 calls to AddPotential, received %d.",
			len(delays), len(fake.Events))
	}
	expected_time := now.Add(1 * time.Millisecond)
	if fake.Events[0].Time != expected_time {
		t.Errorf("Expected first call to be at %s, but was at %s",
			expected_time, fake.Events[0].Time)
	}
	expected_time = now.Add(5 * time.Millisecond)
	if fake.Events[1].Time != expected_time {
		t.Errorf("Expected second call to be at %s, but was at %s",
			expected_time, fake.Events[1].Time)
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
		as <- ActivationEvent{now, makeNeuronWithTerminal(fake, delay, nil, nil)}
	}

	as.ProcessUntilEmpty()

	if len(fake.Events) != len(delays) {
		t.Errorf("Expected %d calls to AddPotential, received %d.",
			len(delays), len(fake.Events))
	}
	for i, delay := range delays {
		expected_time := now.Add(delay)
		if fake.Events[i].Time != expected_time {
			t.Errorf("Expected at %s, got at %s.",
				expected_time, fake.Events[i].Time)
		}
	}
}

func TestActivationStreamAccuracy(t *testing.T) {
	// Connect 1000 neurons, each with a different axon delay,
	// all to the one end-point neuron, so that we can accumulate
	// the accuracy of when the signals reach the end-point.
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

	// Fire all 1000 neurons near-simultaneously (note, the difference
	// in actual firing times won't affect the accuracy which will
	// be calculated based on when each neuron was fired.)
	for _, n := range neurons {
		n.AddPotential(0)
	}
	activation_stream.ProcessUntilEmpty()

	expected_accuracy := time.Duration(100) * time.Microsecond
	if accum.AverageDelta*accum.AverageDelta > expected_accuracy*expected_accuracy {
		t.Errorf("Expected an accuracy better than %s, but "+
			"average delta was %s.", expected_accuracy, accum.AverageDelta)
	}
}

func TestActivationStreamDelay(t *testing.T) {
	// If we string 1000 neurons together, with each axon having the
	// same axon_delay, then we expect the signal to reach the
	// final neuron exactly 1000*axon_delay later. We can then check
	// the variance between when the signal was calculated to reach
	// the end neuron, and when it really did arrive.
	activation_stream := make(ActivationStream, 1000)
	axon_delay := time.Duration(100) * time.Microsecond

	// Create the final end neuron with an event recorder in the
	// axon terminal.
	event_recorder := action_potential.NewEventRecorder(
		new(action_potential.Simple))
	always_fire := action_potential.NewAlwaysFirer(new(action_potential.Simple))
	end := makeNeuronWithTerminal(event_recorder, axon_delay, &activation_stream,
		always_fire)
	prev := end

	// String the other 999 neurons together.
	for i := 0; i < 999; i++ {
		always_fire = action_potential.NewAlwaysFirer(new(action_potential.Simple))
		neuron := makeNeuronWithTerminal(prev, axon_delay,
			&activation_stream, always_fire)
		prev = neuron
	}
	start := prev

	// Start off the chain reaction then sleep for the expected
	// duration of the chain of events.
	started_at := time.Now()
	start.AddPotentialAt(0, started_at)
	go activation_stream.Process()
	expected_duration := axon_delay * 1000
	time.Sleep(expected_duration)
	close(activation_stream)

	if len(event_recorder.Events) != 1 {
		t.Errorf("Expected 1 event, got %d.", len(event_recorder.Events))
	}
	event := event_recorder.Events[0]
	duration := event.Time.Sub(started_at)
	if expected_duration != duration {
		t.Errorf("Expected duration was %s, actual was %s.",
			expected_duration, duration)
	}
	variance := event.RealTime.Sub(event.Time)
	percent_of_expected := float64(variance) / float64(expected_duration) * 100
	tolerance := 0.5
	if math.Abs(percent_of_expected) > tolerance {
		t.Errorf("Actual delay was %f%% of expected delay, which is "+
			"greater than the %f%% tolerance.", percent_of_expected, tolerance)
	}
}
