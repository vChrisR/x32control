package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/hypebeast/go-osc/osc"
)

type oscProcessor struct {
	mapping mixerAddressToChStripMap
}

type MeterData struct {
	_              int32 //bigendian encoded length of blob. Since we know what length to expect we'll ignore this. Rest of packet is littleendian encoded so mixing it just complicated
	NumberOfFloats int32
	Channel        [32]float32
	Aux            [8]float32
	FxReturn       [4]struct {
		Left  float32
		Right float32
	}
	Bus    [16]float32
	Matrix [6]float32
}

func (c oscProcessor) resolveChStrip(address []string) (*ChannelStrip, bool) {
	mixerAddr := fmt.Sprintf("/%v/%v", address[1], address[2])
	chStrip, exists := c.mapping[mixerAddr]
	return chStrip, exists
}

func (c oscProcessor) chHandler(msg *osc.Message) {
	addrElements := strings.Split(msg.Address, "/")
	chStrip, exists := c.resolveChStrip(addrElements)
	if exists {
		var topic, element string
		if len(addrElements) == 5 {
			topic = addrElements[3]
			element = addrElements[4]
		} else {
			if len(addrElements) == 4 {
				topic = "mix"
				element = addrElements[3]
			}
		}
		c.applyMessage(chStrip, topic, element, msg.Arguments)
	}
}

func (c oscProcessor) applyMessage(chStrip *ChannelStrip, topic, element string, arguments []interface{}) {
	switch topic {
	case "mix":

		switch element {
		case "fader":
			value := int(arguments[0].(float32) * 100)
			chStrip.SetFaderValue(value)
		case "on":
			muted := (arguments[0].(int32) == 0)
			chStrip.SetMuted(muted)
		}

	case "config":
		if element == "name" {
			name := arguments[0].(string)
			chStrip.SetLabel(name)
		}
	}
}

func (c oscProcessor) meterHandler(msg *osc.Message) {
	var data MeterData
	if err := binary.Read(bytes.NewBuffer(msg.Arguments[0].([]byte)), binary.LittleEndian, &data); err != nil {
		fmt.Printf(err.Error())
		return
	}

	for key, chStrip := range c.mapping {
		go func(k string, c *ChannelStrip) {
			parts := strings.Split(k, "/")
			index, _ := strconv.Atoi(parts[2])
			index-- //zero based array but 1 based ch nr

			var meterLevel float32
			switch parts[1] {
			case "ch":
				meterLevel = data.Channel[index]
			case "auxin":
				meterLevel = data.Aux[index]
			case "bus":
				meterLevel = data.Bus[index]
			case "mtx":
				meterLevel = data.Matrix[index]
			}

			c.SetMeterValue(meterLevel)
		}(key, chStrip)
	}
}
