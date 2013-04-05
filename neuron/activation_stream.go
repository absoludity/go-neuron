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

// processQueue checks the provided ordered list of terminal events
// processing any which are ready, and returning a timer channel
// which will receive a time when the queue should be processed
// next.
func processQueue(queue *OrderedList) <-chan time.Time {
	e := queue.Front()
	for {
		if e == nil {
			return nil
		}
		te := e.Value.(*TerminalEvent)
		duration := te.Time.Sub(time.Now())
		if duration > 0 {
			return time.NewTimer(duration).C
		}
		next := e.Next()
		queue.Remove(e)
		signalAxonTerminals(te.Neuron.Axon, te.Time)
		e = next
	}
	return nil
}

// Process() processing the incoming activation events, by
// ordering them in a queue and then processing the
// queue.
func (as *ActivationStream) Process() {
	queue := OrderedList{*list.New()}
	// A nil timer channel will block initially, until we assign an
	// alarm.
	var timer_ch <-chan time.Time
	_as := *as
	for {
		select {
		case ae, ok := <-_as:
			if ok {
				terminal_event_time := ae.Time.Add(ae.Neuron.Axon.Delay)
				te := TerminalEvent{terminal_event_time, ae.Neuron}
				queue.Insert(&te)
			} else {
				// No more activation events will be received, but we need to
				// finish processing the queued events. By switching to a nil
				// activation stream, it'll block and allow the remaining
				// queue to be processed.
				_as = nil
			}

			timer_ch = processQueue(&queue)
			if _as == nil && timer_ch == nil {
				return
			}

		case <-timer_ch:
			timer_ch = processQueue(&queue)
			if _as == nil && timer_ch == nil {
				return
			}
		}
	}
}
