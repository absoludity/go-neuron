package action_potential

import (
	"time"
)

// An AccuracyAggregator encapsulates an action potential and
// records the average skew
type AccuracyAccumulator struct {
	ActionPotential
	AverageDelta time.Duration
	Count        int64
}

func NewAccuracyAccumulator(ap ActionPotential) *AccuracyAccumulator {
	return &AccuracyAccumulator{ap, 0, 0}
}

func (f *AccuracyAccumulator) AddPotentialAt(p Potential, t time.Time) (Potential, bool) {
	now := time.Now()
	potential, fired := f.ActionPotential.AddPotentialAt(p, t)
	total_skew := int64(f.AverageDelta)*f.Count + int64(t.Sub(now))
	f.Count += 1
	f.AverageDelta = time.Duration(total_skew / f.Count)
	return potential, fired
}
