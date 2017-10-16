package estim

import (
	"bufio"
	"fmt"
	"io"
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

type ET232 struct {
	r *bufio.Reader
	w io.Writer
}

// NewSerialET232 returns an ET232 that will attempt to communicate with the
// device over the specified serial port (e.g. "COM1" on Windows or
// "/dev/ttyUSB0" on *nix).
func NewSerialET232(portName string) (*ET232, error) {
	c := &serial.Config{Name: portName, Baud: 19200, ReadTimeout: time.Second}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	return &ET232{bufio.NewReader(s), s}, nil
}

func checksum(com []uint8) uint8 {
	var sum uint8 = 0
	for _, b := range com {
		if b >= 0x30 && b <= 0x90 {
			sum = sum + b
		}
	}
	return sum
}

// command sends a command to the device, and returns the response. Commands
// have the following format:
// 1. A single byte representing the type of the command (comType).
// 2. One or more bytes of arguments (args). The number of arguments is
//    determined by the command type.
// 3. A checksum.
//
// The device returns a string ending in '\n'. This function will return the
// response string, excluding the final '\n' character.
func (e *ET232) command(comType uint8, args ...uint8) (string, error) {
	ret := []uint8{comType}
	for _, arg := range args {
		ret = append(ret, []uint8(fmt.Sprintf("%02X", arg))...)
	}
	ret = append(ret, []uint8(fmt.Sprintf("%02X\r", checksum(ret)))...)
	glog.Info("%s", ret)
	_, err := e.w.Write(ret)
	if err != nil {
		return "", err
	}

	str, err := e.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return str[:len(str)-1], nil
}

// ReadAddr reads and returns the value at address addr.
func (e *ET232) ReadAddr(addr uint8) (uint8, error) {
	s, err := e.command(et232ReadCommand, addr)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(s, 16, 8)
	return uint8(v), err
}

// WriteAddr writes val to address addr.
func (e *ET232) WriteAddr(addr uint8, val uint8) error {
	_, err := e.command(et232WriteCommand, addr, val)
	return err
}

// Handshake attempts to perform a serial handshake with the device. This must
// be performed every time the device is powercycled. After calling this
// function, the device should be powercycled for this function to succeed.
func (e *ET232) Handshake() error {
	glog.Info("Attempting to perform a serial handshake.")
	for i := 0; i < 100; i++ {
		// The ET232 sends three bytes "\000CC" when it starts. If trasmission
		// succeeds, the device will listen for serial commands.
		str, _ := e.r.ReadString('\n')
		if str == "\000CC" {
			return nil
		}
	}
	return fmt.Errorf("Failed to connect to the ET232.")
}

// Info reads from the device and summarizes the device state in a human-readable string.
func (e *ET232) Info() (string, error) {
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

// AddCmds adds commands for interacting with the ET232 to the passed Shell.
func (e *ET232) AddCmds(s *ishell.Shell) {
	s.AddCmd(&ishell.Cmd{
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

	s.AddCmd(&ishell.Cmd{
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
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			e.Handshake()
			c.ProgressBar().Stop()
		},
	})
}
