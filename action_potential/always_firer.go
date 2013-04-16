package action_potential

import (
	"time"
)

// An AlwaysFirer encapsulates an action potential but fires
// on every call to AddPotentialAt().
type AlwaysFirer struct {
	ActionPotential
}

func NewAlwaysFirer(ap ActionPotential) *AlwaysFirer {
	return &AlwaysFirer{ap}
}

func (f *AlwaysFirer) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	potential, _ := f.ActionPotential.AddPotentialAt(p, t)
	return potential, true
}

func (f *AlwaysFirer) AddPotential(p Potential) (Potential, bool) {
	return f.AddPotentialAt(p, time.Now())
}
