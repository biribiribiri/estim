// Very simple buttplug server for the ET232.
//
// The "A" and "B" knobs are exposed as vibration motors.
//
// Based on github.com/funjack/golibbuttplug/blob/master/buttplugtest/buttplugtest.go
//
// Example usage:
//   1. go run main.go --port="/dev/ttyUSB0"
//   2. Turn on the ET232.
//   3. Navigate your browser to https://playground.buttplug.world/
//   4. Connect to "ws://localhost:8080/buttplug".
//   5. Enjoy.
//
// To get debugging output, set --logtostderr, and --v=1.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/biribiribiri/buttplugmsg"
	"github.com/biribiribiri/estim"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var handshake = flag.Bool("handshake", true, "perform handshake on start")
var port = flag.String("port", "/dev/ttyUSB0", "serial port name")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// lol security
		return true
	},
}

type ButtplugServer struct {
	Devices []buttplugmsg.Device
}

// Conn is an established websocket connection with the server.
type Conn struct {
	conn   *websocket.Conn
	server *ButtplugServer
}

func (t *ButtplugServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Infof("upgrade error: %v", err)
		return
	}
	defer conn.Close()
	c := &Conn{
		conn:   conn,
		server: t,
	}
	err = c.ReadMessages()
	if err != nil {
		glog.Warningf("error: %v", err)
		return
	}
	return
}

// ReadMessages will read the messages from the websocket to be read and handled.
func (c *Conn) ReadMessages() error {
	for {
		var msgs buttplugmsg.OutgoingMessages
		err := c.conn.ReadJSON(&msgs)
		if _, ok := err.(*websocket.CloseError); ok {
			return err
		} else if err != nil {
			glog.Warningf("error reading buttplugmsg: %v", err)
			continue
		}
		for _, msg := range msgs {
			c.handleMessage(msg)
		}
	}
}

func (c *Conn) handleMessage(m buttplugmsg.OutgoingMessage) {
	switch true {
	case m.RequestServerInfo != nil:
		id := m.RequestServerInfo.ID
		glog.V(1).Infof("<-RequestServerInfo (%d)", id)
		c.sendServerInfo(id)
	case m.RequestDeviceList != nil:
		id := m.RequestDeviceList.ID
		glog.V(1).Infof("<-RequestDeviceList (%d)", id)
		c.sendDeviceList(id)
	case m.StopScanning != nil:
		id := m.StopScanning.ID
		glog.V(1).Infof("<-StopScanning (%d)", id)
		c.sendOk(id)
	case m.Ping != nil:
		id := m.Ping.ID
		glog.V(2).Infof("<-Ping (%d)", id)
		c.sendOk(id)
	case m.VibrateCmd != nil:
		id := m.VibrateCmd.ID
		glog.Info("<-Vibrate ", m.VibrateCmd)
		for _, speed := range m.VibrateCmd.Speeds {
			if speed.Index < uint32(len(knobs)) {
				knobs[speed.Index].Set(speed.Speed)
			}
		}

		c.sendOk(id)
	case m.FleshlightLaunchFW12Cmd != nil:
		id := m.FleshlightLaunchFW12Cmd.ID
		pos, spd := m.FleshlightLaunchFW12Cmd.Position, m.FleshlightLaunchFW12Cmd.Speed
		glog.V(1).Infof("<-FleshlightLaunchFW12Cmd (%d) Postion = %d, Speed = %d", id, pos, spd)
		c.sendOk(id)
	case m.KiirooCmd != nil:
		id := m.KiirooCmd.ID
		glog.V(1).Infof("<-KiirooCmd (%d)", id)
		c.sendOk(id)
	case m.LovenseCmd != nil:
		id := m.LovenseCmd.ID
		glog.V(1).Infof("<-LovenseCmd (%d)", id)
		c.sendOk(id)
	case m.VorzeA10CycloneCmd != nil:
		id := m.VorzeA10CycloneCmd.ID
		spd := m.VorzeA10CycloneCmd.Speed
		clockwise := m.VorzeA10CycloneCmd.Clockwise
		glog.V(1).Infof("<-VorzeA10CycloneCmd (%d) Speed = %d, Clockwise: = %t", id, spd, clockwise)
		c.sendOk(id)
	case m.RawCmd != nil:
		id := m.RawCmd.ID
		glog.V(1).Infof("<-RawCmd (%d)", id)
		c.sendOk(id)
	case m.StartScanning != nil:
		id := m.StartScanning.ID
		glog.V(1).Infof("<-StartScanning, (%d)", id)
		c.sendOk(id)
	case m.StopAllDevices != nil:
		id := m.StopAllDevices.ID
		glog.V(1).Infof("<-StopAllDevices (%d)", id)
		c.sendOk(id)
	case m.StopDeviceCmd != nil:
		id := m.StopDeviceCmd.ID
		glog.V(1).Infof("<-StopDeviceCmd (%d)", id)
		c.sendOk(id)
	}
}

func (c *Conn) sendOk(id uint32) {
	msg := buttplugmsg.IncomingMessage{
		Ok: &buttplugmsg.Empty{
			ID: id,
		},
	}
	err := c.conn.WriteJSON(buttplugmsg.IncomingMessages{msg})
	if err != nil {
		glog.Warningf("error writing: %v", err)
	}
	glog.V(2).Infof("->Ok (%d)", id)
}

func (c *Conn) sendServerInfo(id uint32) {
	msg := buttplugmsg.IncomingMessage{
		ServerInfo: &buttplugmsg.ServerInfo{
			ID:             id,
			ServerName:     "TestButtplug",
			MessageVersion: 1,
			MajorVersion:   1,
			MinorVersion:   0,
			BuildVersion:   0,
			MaxPingTime:    100,
		},
	}
	err := c.conn.WriteJSON(buttplugmsg.IncomingMessages{msg})
	if err != nil {
		glog.Warningf("error writing: %v", err)
	}
	glog.V(1).Infof("->ServerInfo (%d)", id)
}

func (c *Conn) sendDeviceList(id uint32) {
	msg := buttplugmsg.IncomingMessage{
		DeviceList: &buttplugmsg.DeviceList{
			ID:      id,
			Devices: c.server.Devices,
		},
	}
	err := c.conn.WriteJSON(buttplugmsg.IncomingMessages{msg})
	if err != nil {
		glog.Warningf("error writing: %v", err)
	}
	glog.V(1).Infof("->DeviceList (%d)", id)
}

func FatalIfError(err error) {
	if err != nil {
		glog.Fatal(err)
	}
}

var knobs []estim.Knob

func main() {
	flag.Parse()

	e, err := estim.NewSerialET232(*port)
	if err != nil {
		glog.Fatal(err)
	}

	if *handshake {
		log.Println("Performing handshake...")
		err = e.Handshake()
		if err != nil {
			glog.Fatal(err)
		}
	}

	knobs = append(knobs, e.NewKnob(estim.PotA))
	knobs = append(knobs, e.NewKnob(estim.PotB))

	et232Server := &ButtplugServer{
		Devices: []buttplugmsg.Device{
			{
				DeviceName:  "ET232",
				DeviceIndex: 0,
				DeviceMessages: map[string]buttplugmsg.MessageProperty{
					"VibrateCmd":    {FeatureCount: uint32(len(knobs))},
					"StopDeviceCmd": {},
				},
			},
		},
	}

	e.WriteSetting(estim.AnalogOverride, estim.OverrideAB)

	log.Printf("Starting server at: ws://%s/buttplug", *addr)
	http.HandleFunc("/buttplug", et232Server.ServeHTTP)
	glog.Fatal(http.ListenAndServe(*addr, nil))
}
