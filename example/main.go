package main

import (
	"flag"
	"log"
	"time"

	"github.com/biribiribiri/estim"
)

func main() {
	flag.Parse()
	e, err := estim.NewSerialET232("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Performing handshake...")
	e.Handshake() // Perform serial handshake with the device.

	// Force the mode to Intense. Turning off the device while mode is
	// overwritten requires disconnecting power.
	e.WriteSetting(estim.ModeOverride, estim.ModeIntense)

	// Override the A, B, and MA dials.
	e.WriteSetting(estim.AnalogOverride, estim.OverrideAll)

	// Set A and B to 80 (out of 255).
	e.Write(estim.PotA, 80)
	e.Write(estim.PotB, 80)

	time.Sleep(10 * time.Second)

	// Disable mode override so you can turn the thing off.
	e.WriteSetting(estim.ModeOverride, estim.OverrideOff)
}
