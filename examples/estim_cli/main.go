// estim_cli is an interative shell application for interacting with the ET232.
package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/biribiribiri/estim"
)

func addEt232Cmds(e *estim.ET232, s *ishell.Shell) {
	s.AddCmd(&ishell.Cmd{
		Name: "read",
		Help: "read memory address. Ex: read 0x11, read Mode",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("expected 1 argument")
				return
			}
			addr, err := estim.ET232MemString(c.Args[0])
			if err != nil {
				u, err := strconv.ParseUint(c.Args[0], 0, 8)
				if err != nil {
					c.Println(err)
					return
				}
				addr = estim.ET232Mem(u)
			}
			val, err := e.Read(estim.ET232Mem(addr))
			if err != nil {
				c.Println(err)
				return
			}
			c.Printf("0x%02X\n", val)
		},
	})

	s.AddCmd(&ishell.Cmd{
		Name: "write",
		Help: "write memory address. Ex: write 0x1 0x88, write ModeOverride ModeRamp",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("expected 2 arguments")
				return
			}
			// Try to parse a named memory first. If that fails, try parsing
			// as a number.
			addr, err := estim.ET232MemString(c.Args[0])
			if err != nil {
				u, err := strconv.ParseUint(c.Args[0], 0, 8)
				if err != nil {
					c.Println(err)
					return
				}
				addr = estim.ET232Mem(u)
			}

			// Try parsing the value as a named setting. If that doesn't work,
			// try parsing as a number.
			err = e.WriteSetting(addr, estim.ET232Setting(c.Args[1]))
			if err == nil {
				return
			}
			val, err := strconv.ParseUint(c.Args[1], 0, 8)
			if err != nil {
				c.Println(err)
				return
			}
			err = e.Write(estim.ET232Mem(addr), uint8(val))
			if err != nil {
				c.Println(err)
				return
			}
		},
	})

	s.AddCmd(&ishell.Cmd{
		Name: "info",
		Help: "displays info about current device settings",
		Func: func(c *ishell.Context) {
			info, err := e.Info()
			if err != nil {
				c.Println(err)
				return
			}
			c.Println(info)
		},
	})

	s.AddCmd(&ishell.Cmd{
		Name: "handshake",
		Help: "perform a handshake with the device",
		Func: func(c *ishell.Context) {
			c.Println("Performing serial handshake. Please reset the device.")
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			e.Handshake()
			c.ProgressBar().Stop()
		},
	})
}

func main() {
	handshake := flag.Bool("handshake", true, "perform handshake on start")
	port := flag.String("port", "", "serial port name")
	cmds := flag.String("cmds", "", "run in non-interative mode by specifing a list of colon-seperated commands")
	flag.Parse()

	if *port == "" {
		log.Fatal("serial port name must be specified with -port flag")
	}
	e, err := estim.NewSerialET232(*port)
	if err != nil {
		log.Fatal(err)
	}

	shell := ishell.New()
	addEt232Cmds(e, shell)
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
