package action_potential

import (
	"testing"
	"time"
)

func TestAlwaysFirerAddPotentialAt(t *testing.T) {
	potentials := []Potential{0.0, -1.5, 0.1}
	now := time.Now()

	for _, p := range potentials {
		ap := NewAlwaysFirer(new(Simple))

		result, fired := ap.AddPotentialAt(p, now)

		if result != p {
			t.Errorf("Expected potential %.1f, got %.1f.", p, result)
		}
		if !fired {
			t.Error("Expected action potential to fire, but didn't.")
		}

	}
}

func TestAlwaysFirerAddPotential(t *testing.T) {
	potentials := []Potential{0.0, -1.5, 0.1}

	for _, p := range potentials {
		ap := NewAlwaysFirer(new(Simple))

		result, fired := ap.AddPotential(p)

		if result != p {
			t.Errorf("Expected potential %.1f, got %.1f.", p, result)
		}
		if !fired {
			t.Error("Expected action potential to fire, but didn't.")
		}

	}
}
