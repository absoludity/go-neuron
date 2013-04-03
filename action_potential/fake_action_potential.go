package action_potential

import (
	"time"
)

type AddPotentialEvent struct {
	Potential Potential
	Time      time.Time
}

type FakeActionPotential []AddPotentialEvent

func NewFakeActionPotential() FakeActionPotential {
	return make(FakeActionPotential, 0, 10)
}

func (f *FakeActionPotential) GetPotentialAt(now time.Time) Potential {
	return 0.0
}

func (f *FakeActionPotential) GetPotential() Potential {
	return 0.0
}

func (f *FakeActionPotential) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	*f = append(*f, AddPotentialEvent{p, t})
	return 0.0, false
}

func (f *FakeActionPotential) AddPotential(p Potential) (Potential, bool) {
	return f.AddPotentialAt(p, time.Now())
}
