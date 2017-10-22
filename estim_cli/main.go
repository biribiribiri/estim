package main

import (
	"flag"
	"log"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/biribiribiri/estim"
)

func main() {
	handshake := flag.Bool("handshake", true, "perform handshake on start")
	port := flag.String("port", "", "serial port name")
	cmds := flag.String("cmds", "", "a list of colon-seperated commands to pass the estim CLI")
	flag.Parse()

	if *port == "" {
		log.Fatal("serial port name must be specified with -port flag")
	}
	e, err := estim.NewSerialET232(*port)
	if err != nil {
		log.Fatal(err)
	}

	shell := ishell.New()
	e.AddCmds(shell)
	if *handshake {
		err = shell.Process("handshake")
		if err != nil {
			log.Fatal(err)
		}
	}
	if *cmds != "" {
		for _, cmd := range strings.Split(*cmds, ";") {
			s := strings.TrimSpace(cmd)
			shell.Println(s)
			shell.Process(strings.Split(s, " ")...)
		}
	} else {
		shell.Println("estim CLI by biribiribiri. Type \"help\" to get a list of commands.")
		shell.Run()
	}
}
