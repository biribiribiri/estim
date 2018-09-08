package estim

import (
	"fmt"
	"math"
	"time"

	"github.com/golang/glog"
)

type Knob interface {
	Set(val float64) error
	Resolution() float64
}

type KnobQueue struct {
	Knob
	eventQueuer
}

func NewKnobQueue(k Knob) KnobQueue {
	return KnobQueue{k, newEventQueue()}
}

func (k *KnobQueue) Callback(cb func()) {
	k.add(event{0, cb})
}

func (k *KnobQueue) Pulse(val float64, duration time.Duration) {
	glog.V(1).Info("Pulse of ", val, "for ", duration)
	k.add(event{duration, func() {
		if err := k.Set(val); err != nil {
			fmt.Errorf("Set failed: ", err)
		}
	}})
}

func (k *KnobQueue) Ramp(start float64, end float64, duration time.Duration) {
	steps := int(math.Abs(end-start) / k.Resolution())
	if steps == 0 {
		return
	}
	cur := start
	for i := 0; i < steps; i++ {
		k.Pulse(cur, duration/time.Duration(steps))
		cur += math.Copysign(k.Resolution(), end-start)
	}
}

func (k *KnobQueue) WaitDone() {
	k.waitDone()
}

func (k *KnobQueue) Clear() {
	k.clear()
}
