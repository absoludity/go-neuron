package neuron

import (
	"fmt"
	"github.com/absoludity/go-neuron/action_potential"
	"testing"
	"time"
)

// Connect 1000 neurons all to the same action potential, which
// accumulates the accuracy of the timings. Each neuron will
// have a different axon propagation delay.
func BenchmarkActivationStream1000(b *testing.B) {
	b.StopTimer()
	accum := action_potential.NewAccuracyAccumulator(
		new(action_potential.Simple))
	activation_stream := make(ActivationStream, 10000) // ??
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, n := range neurons {
			n.AddPotential(0)
		}

		activation_stream.ProcessUntilEmpty()
	}

	fmt.Print("Average delta:", accum.AverageDelta)
}
