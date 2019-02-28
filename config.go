package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const numChannelStrips byte = 6

type config struct {
	IPAddress     string               `json:"ipAddress"`
	ChannelStrips []channelStripConfig `json:"channelStrips"`
	RecallButton  recallButtonConfig   `json:"recallButton"`
	Langauge      string               `json"langauge:omitEmpty"`
}

type channelStripConfig struct {
	Index      byte   `json:"index"`
	Enabled    bool   `json:"enabled"`
	OscAddress string `json:"oscAddress"`
}

type recallButtonConfig struct {
	Enabled bool   `json:"enabled"`
	SceneNr byte   `json:"sceneNumber"`
	Label   string `json:"label"`
}

func loadConfig(fileName string) (*x32, mixerAddressToChStripMap, []*ChannelStrip, *SceneRecall, string) {
	configFile, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("Unable to load config.json")
	}
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("Unable to read config.json")
	}

	var conf config
	err = json.Unmarshal(configBytes, &conf)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err.Error())
	}

	if len(conf.ChannelStrips) < int(numChannelStrips) {
		log.Fatalf("channelStrips array in config.json needs %v elements", numChannelStrips)
	}

	//Create mixer
	mixer := NewX32(conf.IPAddress, 10023)

	//Create channel strips
	mixerAddressToChStrip := make(mixerAddressToChStripMap)
	allChStrips := make([]*ChannelStrip, numChannelStrips)

	for i := byte(0); i < numChannelStrips; i++ {
		enabled := conf.ChannelStrips[i].Enabled
		chStrip := NewChannelStrip(nil)
		chStrip.index = conf.ChannelStrips[i].Index
		chStrip.SetEnabled(enabled)
		if enabled {
			oscAddress := conf.ChannelStrips[i].OscAddress
			chStrip.mixerChannel = NewX32Channel(oscAddress, mixer)
			mixerAddressToChStrip[oscAddress] = chStrip
		}
		allChStrips[i] = chStrip
	}

	recallButton := NewSceneRecall(nil)
	enabled := conf.RecallButton.Enabled
	recallButton.SetEnabled(enabled)
	if enabled {
		recallButton.scene = conf.RecallButton.SceneNr
		recallButton.mixer = mixer
		recallButton.SetLabel(conf.RecallButton.Label)
		recallButton.mixerChannels = mixerAddressToChStrip
	}
	return mixer, mixerAddressToChStrip, allChStrips, recallButton, conf.Langauge
}
