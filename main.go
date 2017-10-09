package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/golang/glog"
	"github.com/tarm/serial"
)

type mem struct {
	addr     uint8
	desc     string
	min      uint8
	max      uint8
	settings map[string]uint8
}

var et232Mems = map[string]mem{
	"PulseWidthA":      mem{addr: 0x08, desc: "Channel A Pulse Width"},
	"FreqRecA":         mem{addr: 0x09, desc: "Channel A Pulse Frequency Reciprocal"},
	"PulseAmpA":        mem{addr: 0x0A, desc: "Channel A Pulse Amplitude"},
	"PowerCompA":       mem{addr: 0x0B, desc: "Channel A Power Compensation"},
	"PulsePolarityEnA": mem{addr: 0x0C, desc: "Channel A Pulse Enable Polarity"},
	"PulseWidthB":      mem{addr: 0x0E, desc: "Channel B Pulse Width"},
	"FreqRecB":         mem{addr: 0x0F, desc: "Channel B Pulse Frequency Reciprocal"},
	"PulseAmpB":        mem{addr: 0x10, desc: "Channel B Pulse Amplitude"},
	"PowerCompB":       mem{addr: 0x11, desc: "Channel B Power Compensation"},
	"PulsePolarityEnB": mem{addr: 0x12, desc: "Channel B Pulse Enable Polarity"},
	"B":                mem{addr: 0x88, desc: "Position of Pot B"},
	"MA":               mem{addr: 0x89, desc: "Position of MA Pot"},
	"BatteryVoltage":   mem{addr: 0x8A, desc: "Battery Voltage"},
	"AudioInput":       mem{addr: 0x8B, desc: "Audio Input Level"},
	"A":                mem{addr: 0x8C, desc: "Position of Pot A"},
	"Mode": mem{addr: 0xA2, desc: "Mode Switch Position",
		settings: map[string]uint8{
			"Waves":      0x0B,
			"Intense":    0x0A,
			"Random":     0x0E,
			"AudioSoft":  0x06,
			"AudioLoud":  0x02,
			"AudioWaves": 0x03,
			"User":       0x07,
			"HiFreq":     0x05,
			"Climb":      0x01,
			"Throb":      0x00,
			"Combo":      0x04,
			"Thrust":     0x0C,
			"Thump":      0x08,
			"Ramp":       0x09,
			"Stroke":     0x0D,
			"Off":        0x0F}},
	"ModeOverride": mem{addr: 0xA3, desc: "Mode Switch Override"},
	"AnalogOverride": mem{addr: 0xA4, desc: "Analog Input Override",
		settings: map[string]uint8{
			"OverrideAll":  0x1F,
			"OverrideNone": 0x00}},
	"AutoPowerOffTimer":  mem{addr: 0xD3, desc: "Auto Power Off Timer"},
	"ProgramFadeInTimer": mem{addr: 0xD8, desc: "Program Fade In Timer"},
}

const (
	et232WriteCommand = 'I'
	et232ReadCommand  = 'H'
)

type et232 struct {
	*bufio.Reader
	io.Writer
}

func checksum(com []byte) byte {
	var sum byte = 0
	for _, b := range com {
		if b >= 0x30 && b <= 0x90 {
			sum = sum + b
		}
	}
	return sum
}

func (e *et232) command(comType byte, args ...byte) (string, error) {
	ret := []byte{comType}
	for _, arg := range args {
		ret = append(ret, []byte(fmt.Sprintf("%02X", arg))...)
	}
	ret = append(ret, []byte(fmt.Sprintf("%02X\r", checksum(ret)))...)
	glog.Info("%s", ret)
	_, err := e.Write(ret)
	if err != nil {
		return "", err
	}

	str, err := e.ReadString('\n')
	if err != nil {
		return "", err
	}
	return str[:len(str)-1], nil
}

func (e *et232) ReadAddr(addr byte) (byte, error) {
	s, err := e.command(et232ReadCommand, addr)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(s, 16, 8)
	return byte(v), err
}

func (e *et232) WriteAddr(addr byte, val byte) error {
	_, err := e.command(et232WriteCommand, addr, val)
	return err
}

func (e *et232) waitForHandshake() error {
	glog.Info("waiting to connect")
	for i := 0; i < 100; i++ {
		str, _ := e.ReadString('\n')
		if str == "\000CC" {
			return nil
		}
	}
	return fmt.Errorf("Failed to connect")
}

// Info reads from the device and summarizes the device state in a human-readable string.
func (e *et232) Info() (string, error) {
	var out []string
	for name, mem := range et232Mems {
		val, err := e.ReadAddr(mem.addr)
		if err != nil {
			return "", err
		}
		settingName := ""
		if mem.settings != nil {
			settingName = " (Unknown)"
			for name, setting := range mem.settings {
				if setting == val {
					settingName = fmt.Sprintf(" (%s)", name)
				}
			}
		}
		out = append(out, fmt.Sprintf("%s: 0x%02X%s", name, val, settingName))
	}
	sort.Strings(out)
	return strings.Join(out, "\n"), nil
}

func initShell(e *et232) {
	shell := ishell.New()

	// display welcome info.
	shell.Println("estim CLI")

	shell.AddCmd(&ishell.Cmd{
		Name:     "read",
		Help:     "read memory address(es)",
		LongHelp: "read addr1 [addr2] ...",
		Func: func(c *ishell.Context) {
			for _, arg := range c.Args {
				addr, err := strconv.ParseUint(arg, 0, 8)
				if err != nil {
					c.Println(err)
					return
				}
				val, err := e.ReadAddr(uint8(addr))
				if err != nil {
					c.Println(err)
					return
				}
				c.Printf("0x%02X: 0x%02X\n", addr, val)
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name:     "write",
		Help:     "write memory address",
		LongHelp: "write addr value",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("expected 2 arguments")
			}
			addr, err := strconv.ParseUint(c.Args[0], 0, 8)
			if err != nil {
				c.Println(err)
				return
			}
			val, err := strconv.ParseUint(c.Args[1], 0, 8)
			if err != nil {
				c.Println(err)
				return
			}
			err = e.WriteAddr(uint8(addr), uint8(val))
			if err != nil {
				c.Println(err)
				return
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
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

	shell.AddCmd(&ishell.Cmd{
		Name: "handshake",
		Help: "perform a handshake with the device",
		Func: func(c *ishell.Context) {
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			e.waitForHandshake()
			c.ProgressBar().Stop()
		},
	})

	// run shell
	shell.Run()
}

func main() {
	handshake := flag.Bool("handshake", true, "perform handshake on start")

	flag.Parse()
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 19200, ReadTimeout: time.Second}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	e := et232{bufio.NewReader(s), s}

	if *handshake {
		log.Printf("Waiting for handshake...")
		e.waitForHandshake()
	}
	initShell(&e)
}
