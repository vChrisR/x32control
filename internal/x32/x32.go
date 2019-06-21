package x32

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vchrisr/go-osc"
)

type OscHandler func(*osc.Message)

type X32 struct {
	oscClient    osc.Client
	connected    bool
	lastMessage  time.Time
	oscHandlers  map[string]OscHandler
	currentScene int32
}

func New(client osc.Client) *X32 {
	return &X32{
		oscClient:   client,
		connected:   false,
		oscHandlers: make(map[string]OscHandler),
		lastMessage: time.Now().Add(-4 * time.Second),
	}
}

func (x *X32) connect() error {
	log.Printf("Connecting...\n")
	if x.oscClient == nil {
		return fmt.Errorf("No OSC Client configured")
	}
	return x.oscClient.Connect(8)
}

func (x *X32) Send(m *osc.Message) error {
	if x.oscClient == nil {
		return fmt.Errorf("No OSC Client configured")
	}
	return x.oscClient.Send(m)
}

func (x *X32) Start() error {
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

func (x *X32) Stop() {
	x.UnSubscribe()
	// TODO: stop listener, dispatcher, metering
	if x.oscClient != nil {
		x.oscClient.Disconnect()
	}
}

func (x *X32) RecallScene(sceneNr int) error {
	msgAddr := "/-action/goscene"
	msg := osc.NewMessage(msgAddr, int32(sceneNr))
	err := x.Send(msg)
	if err != nil {
		return err
	}

	// give the board half a second to update itself. Probably way too much time but this works.
	time.Sleep(500 * time.Millisecond)

	// seems that xremote stops after scene recall. This triggers xremote immidiately without waiting for potientially 10 seconds.
	x.Send(osc.NewMessage("/xremote"))

	return nil
}

func (x *X32) GetInfo() error {
	msg := osc.NewMessage("/info")
	return x.Send(msg)
}

func (x *X32) StartXremote() {
	xremoteMsg := osc.NewMessage("/xremote")
	x.Send(xremoteMsg)

	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()

		for _ = range tick.C {
			if x.connected {
				x.Send(xremoteMsg)
			}
		}
	}()
}

func (x *X32) RequestMetering() {
	msg := osc.NewMessage("/batchsubscribe")
	msg.Append("/metering")
	msg.Append("/meters/0")
	msg.Append(int32(0))
	msg.Append(int32(0))
	msg.Append(int32(0))

	x.Send(msg)
}

func (x *X32) StartMetering() {
	x.RequestMetering()

	go func() {
		tick := time.NewTicker(9 * time.Second)
		defer tick.Stop()
		renewMsg := osc.NewMessage("/renew")
		renewMsg.Append("/metering")
		for _ = range tick.C {
			if x.connected {
				x.Send(renewMsg)
			}
		}
	}()
}

func (x *X32) UnSubscribe() {
	x.Send(osc.NewMessage("/unsubscribe"))
}

func (x *X32) listen() <-chan *osc.Message {
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

func (x *X32) TrackConnection(connectionLost, disconnected, connected func()) {
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

					if time.Since(x.lastMessage).Seconds() > 20 && !x.connected {
						go connectionLost()
					}

					if time.Since(x.lastMessage).Seconds() > 3 && x.connected {
						x.connected = false
						go disconnected()
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

func (x *X32) Handle(route string, handler OscHandler) {
	x.oscHandlers[route] = handler
}

func (x *X32) startDispatcher(msgChan <-chan *osc.Message) {
	go func() {
		for msg := range msgChan {
			addrElements := strings.Split(msg.Address, "/")
			if handler, exists := x.oscHandlers[addrElements[1]]; exists {
				handler(msg)
			}
		}
	}()
}
