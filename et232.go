package estim

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/tarm/serial"
)

type ET232 struct {
	r    *bufio.Reader
	w    io.Writer
	lock sync.Mutex
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

	return &ET232{r: bufio.NewReader(s), w: s}, nil
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
	glog.V(2).Infof("Sending command: % X", ret)

	e.lock.Lock()
	defer e.lock.Unlock()

	_, err := e.w.Write(ret)
	if err != nil {
		return "", err
	}

	str, err := e.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	glog.V(2).Infof("Recieved response: % X", str)
	return str[:len(str)-1], nil
}

// Read reads and returns the value at address addr.
func (e *ET232) Read(addr ET232Mem) (uint8, error) {
	s, err := e.command(et232ReadCommand, uint8(addr))
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(s, 16, 8)
	return uint8(v), err
}

// Write writes val to the specified memory.
func (e *ET232) Write(mem ET232Mem, val uint8) error {
	glog.V(2).Info("Writing ", val, " to ", mem)
	_, err := e.command(et232WriteCommand, uint8(mem), val)
	return err
}

func floatToUint8(val float64) uint8 {
	if val < 0 {
		val = 0
	}
	if val > 1 {
		val = 1
	}
	return uint8(val * math.MaxUint8)
}

type et232Knob struct {
	mem ET232Mem
	e   *ET232
}

func (e *ET232) NewKnob(mem ET232Mem) Knob {
	return &et232Knob{mem, e}
}

func (k *et232Knob) Set(val float64) error {
	glog.V(1).Infof("Setting %v to %v", k.mem, val)
	return k.e.Write(k.mem, floatToUint8(val))
}

func (k *et232Knob) Resolution() float64 {
	return 1.0 / math.MaxUint8
}

// WriteSetting sets the specified ET232Mem to the passed setting.
func (e *ET232) WriteSetting(mem ET232Mem, setting ET232Setting) error {
	if val, ok := et232SettingMap[et232MemSetting{mem, setting}]; ok {
		return e.Write(mem, val)
	}
	return fmt.Errorf("No setting %s is not valid for memory %v ", setting, mem)
}

// Handshake attempts to perform a serial handshake with the device. This must
// be performed every time the device is powercycled. After calling this
// function, the device should be powercycled for this function to succeed.
func (e *ET232) Handshake() error {
	e.lock.Lock()
	defer e.lock.Unlock()

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

// HandshakeIfNeeded first tries to see if the device is already accepting
// serial commands. If it is not, then it calls Handshake().
func (e *ET232) HandshakeIfNeeded() error {
	if _, err := e.Read(PotA); err == nil {
		return nil
	}
	return e.Handshake()
}

// Info reads from the device and summarizes the device state in a human-readable string.
func (e *ET232) Info() (string, error) {
	mems := []ET232Mem{PulseWidthA,
		FreqRecA,
		PulseAmpA,
		PowerCompA,
		PulsePolarityEnA,
		PulseWidthB,
		FreqRecB,
		PulseAmpB,
		PowerCompB,
		PulsePolarityEnB,
		PotB,
		PotMA,
		BatteryVoltage,
		AudioInput,
		PotA,
		Mode,
		ModeOverride,
		AnalogOverride,
		AutoPowerOffTimer,
		ProgramFadeInTimer}
	var out []string
	for _, mem := range mems {
		val, err := e.Read(mem)
		if err != nil {
			return "", err
		}
		setting, err := GetSetting(mem, val)
		if err == nil {
			out = append(out, fmt.Sprintf("%s: %s", mem, setting))
		} else {
			out = append(out, fmt.Sprintf("%s: 0x%02X", mem, val))
		}
	}
	return strings.Join(out, "\n"), nil
}
