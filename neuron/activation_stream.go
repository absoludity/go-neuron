// go-neuron - A neuron simulator for Go.
//
// Copyright (c) 2013 - Michael Nelson <absoludity@gmail.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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

// processQueue checks the provided queue of terminal events
// processing any which are ready, and returning a timer channel
// which will receive when the queue should be processed
// next.
func processQueue(queue *OrderedList) <-chan time.Time {
	e := queue.Front()
	now := time.Now()
	// How can the delta vary runtime?
	delta := time.Duration(130) * time.Microsecond
	for {
		if e == nil {
			return nil
		}
		te := e.Value.(*TerminalEvent)
		time_until_next := te.Time.Sub(now)
		if time_until_next > delta {
			return time.NewTimer(time_until_next - delta).C
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
// queue. The function returns after the activation stream
// is closed and the queue is cleared.
func (as *ActivationStream) Process() {
	as.process(false)
}

func (as *ActivationStream) ProcessUntilEmpty() {
	as.process(true)
}

func (as *ActivationStream) process(stop_when_empty bool) {
	queue := OrderedList{*list.New()}
	// A nil timer channel will block initially, until we assign an
	// timer channel.
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
			if timer_ch == nil && (stop_when_empty || _as == nil) {
				return
			}

		case <-timer_ch:
			timer_ch = processQueue(&queue)
			if timer_ch == nil && (stop_when_empty || _as == nil) {
				return
			}
		}
	}
}
