package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vchrisr/go-osc"
)

type oscHandler func(*osc.Message)

type x32 struct {
	ip           string
	port         int
	oscClient    osc.Client
	connected    bool
	lastMessage  time.Time
	oscHandlers  map[string]oscHandler
	currentScene int32
}

func NewX32(ip string, port int) *x32 {
	return &x32{
		ip:          ip,
		port:        port,
		oscClient:   osc.NewClient(ip, port, "", 0),
		connected:   false,
		oscHandlers: make(map[string]oscHandler),
	}
}

func (x *x32) connect() error {
	log.Printf("Connecting to %v...\n", x.ip)
	return x.oscClient.Connect(8)
}

func (x *x32) Disconnect() {
	x.UnSubscribe()
	x.oscClient.Disconnect()
}

func (x *x32) Start() error {
	if err := x.connect(); err != nil {
		return err
	}

	if err := x.GetInfo(); err != nil {
		return err
	}

	x.startDispatcher(x.listen())
	x.StartXremote()
	x.StartMetering()

	return nil
}

func (x *x32) RecallScene(sceneNr int) error {
	msgAddr := "/-action/goscene"
	msg := osc.NewMessage(msgAddr, int32(sceneNr))
	err := x.oscClient.Send(msg)
	if err != nil {
		return err
	}

	// give the board half a second to update itself. Probably way too much time but this works.
	time.Sleep(500 * time.Millisecond)

	// seems that xremote stops after scene recall. This triggers xremote immidiately without waiting for potientially 10 seconds.
	x.oscClient.Send(osc.NewMessage("/xremote"))

	return nil
}

func (x *x32) GetInfo() error {
	msg := osc.NewMessage("/info")
	return x.oscClient.Send(msg)
}

func (x *x32) StartXremote() {
	xremoteMsg := osc.NewMessage("/xremote")
	x.oscClient.Send(xremoteMsg)

	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()

		for _ = range tick.C {
			if x.connected {
				x.oscClient.Send(xremoteMsg)
			}
		}
	}()
}

func (x *x32) RequestMetering() {
	msg := osc.NewMessage("/batchsubscribe")
	msg.Append("/metering")
	msg.Append("/meters/0")
	msg.Append(int32(0))
	msg.Append(int32(0))
	msg.Append(int32(0))

	x.oscClient.Send(msg)
}

func (x *x32) StartMetering() {
	x.RequestMetering()

	go func() {
		tick := time.NewTicker(9 * time.Second)
		defer tick.Stop()
		renewMsg := osc.NewMessage("/renew")
		renewMsg.Append("/metering")
		for _ = range tick.C {
			if x.connected {
				x.oscClient.Send(renewMsg)
			}
		}
	}()
}

func (x *x32) UnSubscribe() {
	x.oscClient.Send(osc.NewMessage("/unsubscribe"))
}

func (x *x32) listen() <-chan *osc.Message {
	stream := make(chan *osc.Message, 16)

	go func() {
		defer close(stream)

		for {
			msg, err := x.oscClient.Receive(500 * time.Millisecond)
			if err != nil {
				continue
			}

			x.lastMessage = time.Now()

			select {
			case stream <- msg:
			default:
				log.Println("Receive buffer overrun, message discarded")
			}
		}
	}()

	return stream
}

func (x *x32) TrackConnection(connLost func(), connected func()) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			for _ = range ticker.C {
				if time.Since(x.lastMessage).Seconds() > 1 {
					err := x.GetInfo()
					if err != nil {
						log.Println(err.Error())
					}

					if time.Since(x.lastMessage).Seconds() > 3 && x.connected {
						x.connected = false
						go connLost()
					}
				} else {
					if !x.connected {
						x.connected = true
						go connected()
					}
				}
			}
		}
	}()
}

func (x *x32) Handle(route string, handler oscHandler) {
	x.oscHandlers[route] = handler
}

func (x *x32) startDispatcher(msgChan <-chan *osc.Message) {
	go func() {
		for msg := range msgChan {
			addrElements := strings.Split(msg.Address, "/")
			if handler, exists := x.oscHandlers[addrElements[1]]; exists {
				handler(msg)
			}
		}
	}()
}

type x32channel struct {
	mixer    *x32
	baseAddr string
}

func NewX32Channel(oscAddress string, mixer *x32) x32channel {
	channel := x32channel{
		mixer:    mixer,
		baseAddr: oscAddress,
	}

	return channel
}

func (x x32channel) getMsgAddress(action string) string {
	var msgAddr string
	if strings.Split(x.baseAddr, "/")[1] == "dca" {
		msgAddr = fmt.Sprintf("%v/%v", x.baseAddr, action)
	} else {
		msgAddr = fmt.Sprintf("%v/mix/%v", x.baseAddr, action)
	}

	return msgAddr
}

func (x x32channel) setMute(muted bool) error {
	var value int32
	value = 0
	if !muted {
		value = 1
	}

	msg := osc.NewMessage(x.getMsgAddress("on"), value)
	return x.mixer.oscClient.Send(msg)
}

func (x x32channel) getMute() error {
	msg := osc.NewMessage(x.getMsgAddress("on"))
	return x.mixer.oscClient.Send(msg)
}

func (x x32channel) getName() error {
	msgAddr := fmt.Sprintf("%v/config/name", x.baseAddr)

	msg := osc.NewMessage(msgAddr)
	return x.mixer.oscClient.Send(msg)
}

func (x x32channel) setFaderPosition(position float32) error {
	msg := osc.NewMessage(x.getMsgAddress("fader"), position)
	return x.mixer.oscClient.Send(msg)
}

func (x x32channel) getFaderPosition() error {
	msg := osc.NewMessage(x.getMsgAddress("fader"))
	return x.mixer.oscClient.Send(msg)
}
