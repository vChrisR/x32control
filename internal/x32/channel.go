package x32

import (
	"fmt"
	"strings"

	osc "github.com/vchrisr/go-osc"
)

type Channel struct {
	mixer    *X32
	baseAddr string
}

func NewChannel(oscAddress string, mixer *X32) Channel {
	channel := Channel{
		mixer:    mixer,
		baseAddr: oscAddress,
	}

	return channel
}

func (x Channel) getMsgAddress(action string) string {
	var msgAddr string
	if strings.Split(x.baseAddr, "/")[1] == "dca" {
		msgAddr = fmt.Sprintf("%v/%v", x.baseAddr, action)
	} else {
		msgAddr = fmt.Sprintf("%v/mix/%v", x.baseAddr, action)
	}

	return msgAddr
}

func (x Channel) SetMute(muted bool) error {
	var value int32
	value = 0
	if !muted {
		value = 1
	}

	msg := osc.NewMessage(x.getMsgAddress("on"), value)
	return x.mixer.oscClient.Send(msg)
}

func (x Channel) GetMute() error {
	msg := osc.NewMessage(x.getMsgAddress("on"))
	return x.mixer.oscClient.Send(msg)
}

func (x Channel) GetName() error {
	msgAddr := fmt.Sprintf("%v/config/name", x.baseAddr)

	msg := osc.NewMessage(msgAddr)
	return x.mixer.oscClient.Send(msg)
}

func (x Channel) SetFaderPosition(position float32) error {
	msg := osc.NewMessage(x.getMsgAddress("fader"), position)
	return x.mixer.oscClient.Send(msg)
}

func (x Channel) GetFaderPosition() error {
	msg := osc.NewMessage(x.getMsgAddress("fader"))
	return x.mixer.oscClient.Send(msg)
}
