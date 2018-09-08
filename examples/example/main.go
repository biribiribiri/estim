package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/biribiribiri/estim"
	"github.com/golang/glog"
)

func main() {
	port := flag.String("port", "/dev/ttyUSB0", "serial port name")

	flag.Parse()
	et232, err := estim.NewSerialET232(*port)
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Println("Performing handshake...")
	// Perform serial handshake with the device if it's not connected already.
	// If it doesn't connect right away, then turn the device off and on
	// again.
	err = et232.HandshakeIfNeeded()
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Println("Connected!")

	// Override the A, B, and MA dials.
	err = et232.WriteSetting(estim.AnalogOverride, estim.OverrideAll)
	if err != nil {
		glog.Fatal(err)
	}

	// Set A and B to 80 (out of 255).
	err = et232.Write(estim.PotA, 80)
	if err != nil {
		glog.Fatal(err)
	}
	err = et232.Write(estim.PotB, 80)
	if err != nil {
		glog.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// A "knob" lets you control a setting by specifying a value between 0.0
	// and 1.0.
	aKnob := et232.NewKnob(estim.PotA)

	// You can create a KnobQueue for controlling the position of a knob over
	// a period of time. Operations on a KnobQueue are non-blocking.
	aKnobQueue := estim.NewKnobQueue(aKnob)

	// Operations will occur in the order that they are added to the queue.
	// Increase A from 20% to 40% over a period of 3 seconds. Then set A to
	// 50% for 1 second. Then log a message, and set A to 0%.
	aKnobQueue.Ramp(0.2, 0.4, 3*time.Second)
	aKnobQueue.Pulse(0.5, time.Second)
	aKnobQueue.Callback(func() { fmt.Println("Setting A to 0!") })
	aKnobQueue.Pulse(0, time.Second)

	// Block until all queued operations are done.
	aKnobQueue.WaitDone()

	// Disable mode override so you can turn the thing off.
	et232.WriteSetting(estim.ModeOverride, estim.OverrideOff)
}
