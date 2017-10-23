package estim

import "fmt"

// ET232Mem represents an ET232 memory address.
type ET232Mem uint8

const (
	PulseWidthA        ET232Mem = 0x08 // PulseWidthA is Channel A Pulse Width
	FreqRecA           ET232Mem = 0x09 // FreqRecA is Channel A Pulse Frequency Reciprocal
	PulseAmpA          ET232Mem = 0x0A // PulseAmpA is Channel A Pulse Amplitude
	PowerCompA         ET232Mem = 0x0B // PowerCompA is Channel A Power Compensation
	PulsePolarityEnA   ET232Mem = 0x0C // PulsePolarityEnA is Channel A Pulse Enable Polarity
	PulseWidthB        ET232Mem = 0x0E // PulseWidthB is Channel B Pulse Width
	FreqRecB           ET232Mem = 0x0F // FreqRecB is Channel B Pulse Frequency Reciprocal
	PulseAmpB          ET232Mem = 0x10 // PulseAmpB is Channel B Pulse Amplitude
	PowerCompB         ET232Mem = 0x11 // PowerCompB is Channel B Power Compensation
	PulsePolarityEnB   ET232Mem = 0x12 // PulsePolarityEnB is Channel B Pulse Enable Polarity
	PotB               ET232Mem = 0x88 // PotB is Position of Pot B
	PotMA              ET232Mem = 0x89 // PotMA is Position of MA Pot
	BatteryVoltage     ET232Mem = 0x8A // BatteryVoltage is Battery Voltage
	AudioInput         ET232Mem = 0x8B // AudioInput is Audio Input Level
	PotA               ET232Mem = 0x8C // PotA is Position of Pot A
	Mode               ET232Mem = 0xA2 // Mode is Mode Switch Position
	ModeOverride       ET232Mem = 0xA3 // ModeOverride is Mode Switch Override
	AnalogOverride     ET232Mem = 0xA4 // AnalogOverride is Analog Input Override
	AutoPowerOffTimer  ET232Mem = 0xD3 // AutoPowerOffTimer is Auto Power Off Timer
	ProgramFadeInTimer ET232Mem = 0xD8 // ProgramFadeInTimer is Program Fade In Timer
)

// ET232Setting is a named setting for an ET232Mem. This is used for certain
// ET232Mems that only take certain discrete values (e.g. the Mode switch).
type ET232Setting string

const (
	ModeWaves      ET232Setting = "ModeWaves"      // Waves setting for "Mode" and "ModeOverride"
	ModeIntense    ET232Setting = "ModeIntense"    // Intense setting for "Mode" and "ModeOverride"
	ModeRandom     ET232Setting = "ModeRandom"     // Random setting for "Mode" and "ModeOverride"
	ModeAudioSoft  ET232Setting = "ModeAudioSoft"  // AudioSoft setting for "Mode" and "ModeOverride"
	ModeAudioLoud  ET232Setting = "ModeAudioLoud"  // AudioLoud setting for "Mode" and "ModeOverride"
	ModeAudioWaves ET232Setting = "ModeAudioWaves" // AudioWaves setting for "Mode" and "ModeOverride"
	ModeUser       ET232Setting = "ModeUser"       // User setting for "Mode" and "ModeOverride"
	ModeHiFreq     ET232Setting = "ModeHiFreq"     // HiFreq setting for "Mode" and "ModeOverride"
	ModeClimb      ET232Setting = "ModeClimb"      // Climb setting for "Mode" and "ModeOverride"
	ModeThrob      ET232Setting = "ModeThrob"      // Throb setting for "Mode" and "ModeOverride"
	ModeCombo      ET232Setting = "ModeCombo"      // Combo setting for "Mode" and "ModeOverride"
	ModeThrust     ET232Setting = "ModeThrust"     // Thrust setting for "Mode" and "ModeOverride"
	ModeThump      ET232Setting = "ModeThump"      // Thump setting for "Mode" and "ModeOverride"
	ModeRamp       ET232Setting = "ModeRamp"       // Ramp setting for "Mode" and "ModeOverride"
	ModeStroke     ET232Setting = "ModeStroke"     // Stroke setting for "Mode" and "ModeOverride"
	ModeOff        ET232Setting = "ModeOff"        // Off setting for "Mode" and "ModeOverride"

	OverrideAll ET232Setting = "OverrideAll"
	OverrideOff ET232Setting = "OverrideOff"
)

//go:generate enumer -type=ET232Mem

type et232MemSetting struct {
	mem     ET232Mem
	setting ET232Setting
}

var et232SettingMap = map[et232MemSetting]uint8{
	et232MemSetting{Mode, ModeWaves}:      0x0B,
	et232MemSetting{Mode, ModeIntense}:    0x0A,
	et232MemSetting{Mode, ModeRandom}:     0x0E,
	et232MemSetting{Mode, ModeAudioSoft}:  0x06,
	et232MemSetting{Mode, ModeAudioLoud}:  0x02,
	et232MemSetting{Mode, ModeAudioWaves}: 0x03,
	et232MemSetting{Mode, ModeUser}:       0x07,
	et232MemSetting{Mode, ModeHiFreq}:     0x05,
	et232MemSetting{Mode, ModeClimb}:      0x01,
	et232MemSetting{Mode, ModeThrob}:      0x00,
	et232MemSetting{Mode, ModeCombo}:      0x04,
	et232MemSetting{Mode, ModeThrust}:     0x0C,
	et232MemSetting{Mode, ModeThump}:      0x08,
	et232MemSetting{Mode, ModeRamp}:       0x09,
	et232MemSetting{Mode, ModeStroke}:     0x0D,
	et232MemSetting{Mode, ModeOff}:        0x0F,

	et232MemSetting{ModeOverride, OverrideOff}:    0x00,
	et232MemSetting{ModeOverride, ModeWaves}:      0x8B,
	et232MemSetting{ModeOverride, ModeIntense}:    0x8A,
	et232MemSetting{ModeOverride, ModeRandom}:     0x8E,
	et232MemSetting{ModeOverride, ModeAudioSoft}:  0x86,
	et232MemSetting{ModeOverride, ModeAudioLoud}:  0x82,
	et232MemSetting{ModeOverride, ModeAudioWaves}: 0x83,
	et232MemSetting{ModeOverride, ModeUser}:       0x87,
	et232MemSetting{ModeOverride, ModeHiFreq}:     0x85,
	et232MemSetting{ModeOverride, ModeClimb}:      0x81,
	et232MemSetting{ModeOverride, ModeThrob}:      0x80,
	et232MemSetting{ModeOverride, ModeCombo}:      0x84,
	et232MemSetting{ModeOverride, ModeThrust}:     0x8C,
	et232MemSetting{ModeOverride, ModeThump}:      0x88,
	et232MemSetting{ModeOverride, ModeRamp}:       0x89,
	et232MemSetting{ModeOverride, ModeStroke}:     0x8D,
	et232MemSetting{ModeOverride, ModeOff}:        0x8F,

	et232MemSetting{AnalogOverride, OverrideAll}: 0x8D,
	et232MemSetting{AnalogOverride, OverrideOff}: 0x8F,
}

// GetSetting returns the ET232Setting corresponding to a specified ET232Mem and value.
func GetSetting(mem ET232Mem, val uint8) (ET232Setting, error) {
	for ms, v := range et232SettingMap {
		if ms.mem == mem && v == val {
			return ms.setting, nil
		}
	}
	return "", fmt.Errorf("No ET232Setting for memory %v and value %v", mem, val)
}

const (
	et232WriteCommand = 'I'
	et232ReadCommand  = 'H'
)
