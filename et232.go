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
		Name: "read",
		Help: "read memory address. Ex: read 0x11",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("expected 1 argument")
				return
			}
			addr, err := ET232MemString(c.Args[0])
			if err != nil {
				u, err := strconv.ParseUint(c.Args[0], 0, 8)
				if err != nil {
					c.Println(err)
					return
				}
				addr = ET232Mem(u)
			}
			val, err := e.ReadAddr(uint8(addr))
			if err != nil {
				c.Println(err)
				return
			}
			c.Printf("0x%02X\n", val)
		},
	})

	s.AddCmd(&ishell.Cmd{
		Name: "write",
		Help: "write memory address. Ex: write 0x1 0x88",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("expected 2 arguments")
				return
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
			c.Println("Performing serial handshake. Please reset the device.")
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			e.Handshake()
			c.ProgressBar().Stop()
		},
	})
}
