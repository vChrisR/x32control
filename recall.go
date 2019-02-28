package main

import (
	"fmt"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/therecipe/qt/core"
)

type SceneRecall struct {
	core.QObject

	client        *osc.Client
	scene         byte
	mixer         *x32
	mixerChannels mixerAddressToChStripMap

	_ bool   `property:"enabled"`
	_ string `property:"label"`
	_ func() `constructor:"init"`
	_ func() `signal:"recallClicked,auto"`
}

func (r *SceneRecall) init() {
	return
}

func (r *SceneRecall) recallClicked() {
	if err := r.mixer.RecallScene(int(r.scene)); err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(500 * time.Millisecond)

	for _, channel := range r.mixerChannels {
		channel.updateFromMixer()
	}
}
