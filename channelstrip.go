package main

import (
	"github.com/therecipe/qt/core"
)

type mixerAddressToChStripMap map[string]*ChannelStrip

type ChannelStrip struct {
	core.QObject

	index             byte
	mixerChannel      x32channel
	cmdsSend          int
	lastFaderPosition float32

	_ bool          `property:"muted"`
	_ int           `property:"faderValue"`
	_ float32       `property:"meterValue"`
	_ string        `property:"label"`
	_ func()        `constructor:"init"`
	_ func(bool)    `slot:"muteclicked,auto"`
	_ func(float32) `slot:"fadermoved,auto"`
}

func (c *ChannelStrip) init() {
	c.SetMuted(true)
	c.SetFaderValue(0)
	c.SetMeterValue(0)
}

func (c *ChannelStrip) updateFromMixer() error {
	err := c.mixerChannel.getMute()
	err = c.mixerChannel.getFaderPosition()
	err = c.mixerChannel.getName()
	return err
}

func (c *ChannelStrip) muteclicked(checked bool) {
	c.SetMuted(checked)
	c.mixerChannel.setMute(checked)
}

func (c *ChannelStrip) fadermoved(value float32) {
	//Round to value x32 can acutally process
	intval := int(value * 1023.5)
	faderPosition := float32(intval) / 1023

	//Some fader moves might produce the same fader position afder rounding to an x32 value. If the previous value is the same as current we won't send anything.
	if faderPosition != c.lastFaderPosition {
		c.mixerChannel.setFaderPosition(faderPosition)
		c.lastFaderPosition = faderPosition
	}
}
