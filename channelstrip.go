package main

import (
	"github.com/vchrisr/x32control/internal/x32"

	"github.com/therecipe/qt/core"
)

type ChannelStrip struct {
	qmlObj            *core.QObject
	mixerChannel      x32.Channel
	lastFaderPosition float32
}

type ChannelStrips map[string]*ChannelStrip

func (c *ChannelStrip) SetFaderValue(value int) {
	c.qmlObj.SetProperty("faderValue", core.NewQVariant7(value))
}

func (c *ChannelStrip) SetMuted(value bool) {
	c.qmlObj.SetProperty("muted", core.NewQVariant11(value))
}

func (c *ChannelStrip) SetLabel(value string) {
	c.qmlObj.SetProperty("label", core.NewQVariant14(value))
}

func (c *ChannelStrip) SetMeterValue(value float32) {
	sendValue := value
	if sendValue < 0.25 && sendValue > 0 { //meters are tiny. show 25% when any signal is detected
		sendValue = 0.25
	}
	c.qmlObj.SetProperty("meterValue", core.NewQVariant13(sendValue))
}

func (c *ChannelStrip) updateFromMixer() error {
	err := c.mixerChannel.GetMute()
	err = c.mixerChannel.GetFaderPosition()
	err = c.mixerChannel.GetName()

	return err
}

func (c *ChannelStrip) sendFaderPosition(value float32) {
	//Round to value x32 can acutally process
	intval := int(value * 1023.5)
	faderPosition := float32(intval) / 1023

	//Some fader moves might produce the same fader position afder rounding to an x32 value. If the previous value is the same as current we won't send anything.
	if faderPosition != c.lastFaderPosition {
		c.mixerChannel.SetFaderPosition(faderPosition)
		c.lastFaderPosition = faderPosition
	}
}
