package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/quickcontrols2"
)

func main() {
	//setup some Qt stuff
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	gui.NewQGuiApplication(len(os.Args), os.Args)

	quickcontrols2.QQuickStyle_SetStyle("Material")
	view := quick.NewQQuickView(nil)
	view.SetTitle("x32control")
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)

	//Load config
	mixer, channelStrips, allChStrips, recallButton, language := loadConfig("config.json")

	//Set langauge
	if language != "" && language != "en" {
		var translator = core.NewQTranslator(nil)
		if loaded := translator.Load(fmt.Sprintf("qml_%v", language), ":/qml", "", ""); loaded == false {
			fmt.Println("unable to load language file for selected language")
		}
		core.QCoreApplication_InstallTranslator(translator)
	}

	//Configure mixer
	chStripProcessor := oscProcessor{
		mapping: channelStrips,
	}

	mixer.Handle("ch", chStripProcessor.chHandler)
	mixer.Handle("main", chStripProcessor.chHandler)
	mixer.Handle("dca", chStripProcessor.chHandler)
	mixer.Handle("bus", chStripProcessor.chHandler)
	mixer.Handle("auxin", chStripProcessor.chHandler)
	mixer.Handle("mtx", chStripProcessor.chHandler)
	mixer.Handle("metering", chStripProcessor.meterHandler)

	//make sure to disconnect when app is closed gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigs
		fmt.Println("Closing...")
		mixer.Disconnect()
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	//create al the context objects for QML
	for _, chStrip := range allChStrips {
		view.RootContext().SetContextProperty(fmt.Sprintf("strip%v", chStrip.index), chStrip)
	}

	view.RootContext().SetContextProperty("recall", recallButton)
	bi := NewBusyIndicator(nil)
	view.RootContext().SetContextProperty("busy", bi)
	status := NewControllerStatus(nil)
	view.RootContext().SetContextProperty("controllerStatus", status)

	//load  the qml
	view.SetSource(core.NewQUrl3("qrc:/qml/main.qml", 0))
	view.Show()

	//start the mixer interaction
	if err := mixer.Start(); err != nil {
		panic(err.Error())
	}

	//Track the mixer connection. Functions for onDisconnect and onConnect are passed
	mixer.TrackConnection(
		func() {
			fmt.Println("Disconnected")
			bi.SetBusy(true)
			for _, chStrip := range channelStrips {
				chStrip.lastFaderPosition = 0
			}
		},
		func() {
			fmt.Println("Connected")
			for _, chStrip := range channelStrips {
				chStrip.updateFromMixer()
			}

			bi.SetBusy(false)

			mixer.Send(osc.NewMessage("/xremote"))
			mixer.RequestMetering()
		})

	gui.QGuiApplication_Exec()
}
