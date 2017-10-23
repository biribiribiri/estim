package estim

type mem struct {
	addr uint8
	desc string

	// If not nil, this memory only takes discrete values.
	settings map[string]uint8
}

// ET232Mem represents an ET232 memory address.
type ET232Mem uint8

const (
	PulseWidthA        ET232Mem = 0x08 // Channel A Pulse Width
	FreqRecA           ET232Mem = 0x09 // Channel A Pulse Frequency Reciprocal
	PulseAmpA          ET232Mem = 0x0A // Channel A Pulse Amplitude
	PowerCompA         ET232Mem = 0x0B // Channel A Power Compensation
	PulsePolarityEnA   ET232Mem = 0x0C // Channel A Pulse Enable Polarity
	PulseWidthB        ET232Mem = 0x0E // Channel B Pulse Width
	FreqRecB           ET232Mem = 0x0F // Channel B Pulse Frequency Reciprocal
	PulseAmpB          ET232Mem = 0x10 // Channel B Pulse Amplitude
	PowerCompB         ET232Mem = 0x11 // Channel B Power Compensation
	PulsePolarityEnB   ET232Mem = 0x12 // Channel B Pulse Enable Polarity
	PotB               ET232Mem = 0x88 // Position of Pot B
	PotMA              ET232Mem = 0x89 // Position of MA Pot
	BatteryVoltage     ET232Mem = 0x8A // Battery Voltage
	AudioInput         ET232Mem = 0x8B // Audio Input Level
	PotA               ET232Mem = 0x8C // Position of Pot A
	Mode               ET232Mem = 0xA2 // Mode Switch Position
	ModeOverride       ET232Mem = 0xA3 // Mode Switch Override
	AnalogOverride     ET232Mem = 0xA4 // Analog Input Override
	AutoPowerOffTimer  ET232Mem = 0xD3 // Auto Power Off Timer
	ProgramFadeInTimer ET232Mem = 0xD8 // Program Fade In Timer
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
			"OverrideAll": 0x1F,
			"OverrideOff": 0x00}},
	"AutoPowerOffTimer":  mem{addr: 0xD3, desc: "Auto Power Off Timer"},
	"ProgramFadeInTimer": mem{addr: 0xD8, desc: "Program Fade In Timer"},
}

const (
	et232WriteCommand = 'I'
	et232ReadCommand  = 'H'
)
