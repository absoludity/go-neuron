package action_potential

import (
	"time"
)

type AddPotentialEvent struct {
	Potential Potential
	Time      time.Time
}

type FakeEventRecorder struct {
	Events []AddPotentialEvent
}

func (f *FakeEventRecorder) GetPotentialAt(now time.Time) Potential {
	return 0.0
}

func (f *FakeEventRecorder) GetPotential() Potential {
	return 0.0
}

func (f *FakeEventRecorder) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	if f.Events == nil {
		f.Events = make([]AddPotentialEvent, 0, 10)
	}
	f.Events = append(f.Events, AddPotentialEvent{p, t})
	return 0.0, false
}

func (f *FakeEventRecorder) AddPotential(p Potential) (Potential, bool) {
	return f.AddPotentialAt(p, time.Now())
}
