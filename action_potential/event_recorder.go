package action_potential

import (
	"time"
)

type AddPotentialEvent struct {
	Potential Potential
	Time      time.Time
}

// An EventRecorder encapsulates an action potential and records
// each AddPotential event that occurs.
type EventRecorder struct {
	ActionPotential
	Events []AddPotentialEvent
}

func NewEventRecorder(ap ActionPotential) *EventRecorder {
	return &EventRecorder{ap, make([]AddPotentialEvent, 0, 10)}
}

func (f *EventRecorder) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	f.Events = append(f.Events, AddPotentialEvent{p, t})
	return f.ActionPotential.AddPotentialAt(p, t)
}