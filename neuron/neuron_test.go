package neuron

import (
	"github.com/absoludity/go-neuron/action_potential"
	"testing"
	"time"
)

func TestNeuronFire(t *testing.T) {
	now := time.Now()
	cb := action_potential.New(1, now)
	as := make(ActivationStream, 1)
	n := &Neuron{Axon{}, as, cb}
	at := now.Add(time.Microsecond * 5)

	actual_potential := n.AddPotentialAt(action_potential.THRESHOLD_POTENTIAL, at)

	if actual_potential != action_potential.PEAK_POTENTIAL {
		t.Error("Expected potential:", action_potential.PEAK_POTENTIAL,
			"Actual potential: ", actual_potential)
	}
	ae, ok := <-as
	if !ok {
		t.Errorf("No activation event received on activation stream.")
	}
	if n != ae.Neuron {
		t.Errorf("Expected activation event for %s, but received %s",
			n, ae.Neuron)
	}
	if at != ae.Time {
		t.Errorf("Received activation event for incorrect time.")
	}
}

func TestNeuronNoFire(t *testing.T) {
	now := time.Now()
	cb := action_potential.New(1, now)
	as := make(ActivationStream, 1)
	n := &Neuron{Axon{}, as, cb}
	at := now.Add(time.Microsecond * 5)

	actual_potential := n.AddPotentialAt(5, at)

	if actual_potential != 6 {
		t.Error("Expected potential:", 6,
			"Actual potential: ", actual_potential)
	}

	// Close the activation stream so it doesn't block when empty.
	close(as)
	ae := <-as
	nil_activation_event := ActivationEvent{}
	if ae != nil_activation_event {
		t.Error("Expected the activation stream to be empty, but found ",
			ae)
	}
}
