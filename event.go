package estim

import (
	"sync"
	"time"
)

type event struct {
	duration time.Duration
	action   func()
}

type eventQueue struct {
	queue     []event
	lock      sync.Mutex
	done      *sync.Cond
	clearChan chan bool
}

type eventQueuer interface {
	add(e event)
	clear()
	waitDone()
}

func newEventQueue() eventQueuer {
	eq := eventQueue{clearChan: make(chan bool)}
	eq.done = sync.NewCond(&eq.lock)
	process := func() time.Duration {
		eq.lock.Lock()
		defer eq.lock.Unlock()
		if len(eq.queue) == 0 {
			eq.done.Broadcast()
			return 10 * time.Millisecond
		}
		var e event
		e, eq.queue = eq.queue[0], eq.queue[1:]
		e.action()
		return e.duration
	}

	go func() {
		for {
			select {
			case <-time.After(process()):
			case <-eq.clearChan:
			}
		}
	}()
	return &eq
}

func (eq *eventQueue) add(e event) {
	eq.lock.Lock()
	defer eq.lock.Unlock()

	eq.queue = append(eq.queue, e)
}

func (eq *eventQueue) clear() {
	eq.lock.Lock()
	defer eq.lock.Unlock()

	eq.queue = nil
	eq.clearChan <- true
}

func (eq *eventQueue) waitDone() {
	eq.lock.Lock()
	eq.done.Wait()
}
