package main

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	IPAddress     string               `json:"ipAddress"`
	ChannelStrips []channelStripConfig `json:"channelStrips"`
	RecallButton  recallButtonConfig   `json:"recallButton"`
	Language      string               `json:"language:omitEmpty"`
}

type channelStripConfig struct {
	OscAddress string `json:"oscAddress"`
}

type recallButtonConfig struct {
	Enabled bool   `json:"enabled"`
	SceneNr byte   `json:"sceneNumber"`
	Label   string `json:"label"`
}

const maxNumChannelStrips = 9

func loadConfig(fileName string) (*x32, config) {
	configFile, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("Unable to load config.json")
	}
	defer configFile.Close()

	var conf config
	err = json.NewDecoder(configFile).Decode(&conf)
	if err != nil {
		log.Fatalf("Unable to read config.json: %v", err.Error())
	}

	if len(conf.ChannelStrips) == 0 {
		log.Fatalf("No channelStrips configured in config.json.")
	}

	if len(conf.ChannelStrips) > int(maxNumChannelStrips) {
		log.Printf("channelStrips array in config.json contains too many elements. Ony %v channel strips supported. Remaining strips will be ignored.\n", maxNumChannelStrips)
		conf.ChannelStrips = conf.ChannelStrips[:maxNumChannelStrips]
	}

	//Create mixer
	mixer := NewX32(conf.IPAddress, 10023)

	return mixer, conf
}
