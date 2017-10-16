package main

import (
	"flag"
	"log"

	"github.com/abiosoft/ishell"
	"github.com/biribiribiri/estim"
)

func main() {
	handshake := flag.Bool("handshake", true, "perform handshake on start")

	flag.Parse()
	e, err := estim.NewSerialET232("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err)
	}

	if *handshake {
		log.Printf("Performing serial handshake. Please reset the device.")
		e.Handshake()
	}

	shell := ishell.New()
	shell.Println("estim CLI")
	e.AddCmds(shell)
	shell.Run()
}
