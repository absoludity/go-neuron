package neuron

import (
	"container/list"
	"time"
)

// An ActivationEvent records the neuron and time at which it
// was activated.
type ActivationEvent struct {
	Time   time.Time
	Neuron *Neuron
}

type OrderedList struct {
	list.List
}

func (l *OrderedList) Insert(value interface{}) *list.Element {
	// Do benchmarks with built-in sort algorithm too. Very special
	// case insert into sorted list.
	event_time := value.(*TerminalEvent).Time
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*TerminalEvent).Time.After(event_time) {
			return l.InsertBefore(value, e)
		}
	}
	return l.PushBack(value)
}

// A TerminalEvent records the neuron and the time at which
// the signal reaches the axon terminals.
type TerminalEvent ActivationEvent

// An ActivationStream communicates the activation events for further
// processing.
type ActivationStream chan ActivationEvent

func signalAxonTerminals(a Axon, t time.Time) {
	for _, n := range a.Terminals {
		// Should the potential for each be relative to total
		// potential, or constant, or divided by the num of terminals?
		n.AddPotentialAt(5.0, t)
	}
}

// checkQueue checks the provided ordered list of terminal events
// processing any which are ready, and returning the time when the
// next check should occur.
func checkQueue(queue *OrderedList) time.Duration {
	for e := queue.Front(); e != nil; e = e.Next() {
		te := e.Value.(*TerminalEvent)
		duration := te.Time.Sub(time.Now())
		if duration > 0 {
			return duration
		}
		queue.Remove(e)
		signalAxonTerminals(te.Neuron.Axon, te.Time)
	}
	return 0
}

func (as *ActivationStream) Process() {
	queue := OrderedList{*list.New()}
	time_ch := make(<-chan time.Time)
	var ae ActivationEvent
	var ok bool
	var next_check_in time.Duration
	for {
		select {
		case ae, ok = <-*as:
			if !ok {
				// No more events will be received, but we need
				// to finish processing the queued events.
				for {
					if queue.Len() == 0 {
						return
					} else {
						// Don't really need channel here, as no select
						// required, switch to simple timer?
						time_ch = time.NewTimer(next_check_in).C
					}
					<-time_ch
					next_check_in = checkQueue(&queue)
				}
			}
			terminal_event_time := ae.Time.Add(ae.Neuron.Axon.Delay)
			te := TerminalEvent{terminal_event_time, ae.Neuron}
			queue.Insert(&te)

			next_check_in = checkQueue(&queue)
			if next_check_in > 0 {
				time_ch = time.NewTimer(next_check_in).C
			}

		case <-time_ch:
			next_check_in := checkQueue(&queue)
			if next_check_in > 0 {
				time_ch = time.NewTimer(next_check_in).C
			}
		}
	}
}
