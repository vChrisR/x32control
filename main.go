package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/quickcontrols2"
	osc "github.com/vchrisr/go-osc"
	"github.com/vchrisr/x32control/internal/x32"
	//	_ "net/http/pprof"
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
	conf := loadConfig("config.json")

	//Set langauge
	if conf.Language != "" && conf.Language != "en" {
		var translator = core.NewQTranslator(nil)
		if loaded := translator.Load(fmt.Sprintf("qml_%v", conf.Language), ":/qml", "", ""); loaded == false {
			log.Println("unable to load language file for selected language")
		}
		core.QCoreApplication_InstallTranslator(translator)
	}

	// Auto discover mixers on the networks
	var mixer *x32.X32
	if conf.IPAddress == "" || conf.IPAddress == "auto" {
		log.Println("No x32 IP configured. Performing AutoDiscover")
		discoveredIps, err := x32.AutoDiscover(60) //try AD for 60 seconds
		if err != nil {
			log.Printf("Error while running AutoDiscover: %v Exitting now.", err) //if nothing found after 60secs just exit.
			os.Exit(1)
		}

		if len(discoveredIps) > 1 {
			log.Fatalf("More than one x32 discovered. This is currently not supported. Please configure a mixer ip in config.json")
		}

		log.Printf("Discovered x32: %v", discoveredIps[0])
		conf.IPAddress = strings.Split(discoveredIps[0], ":")[0]
	} else {
		log.Printf("IP configured: %v. Skipping AutoDiscover", conf.IPAddress)
	}

	client := osc.NewClient(conf.IPAddress, 10023, "", 0)
	mixer = x32.New(client)

	qmlRoot := initQmlRoot(view, conf, mixer)

	//Configure mixer
	chStripProcessor := NewOscStripProcessor(qmlRoot)
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
		log.Println("Closing...")
		mixer.Stop()
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	//load  the qml
	view.SetSource(core.NewQUrl3("qrc:/qml/main.qml", 0))
	view.Show()

	//start the mixer interaction
	if err := mixer.Start(); err != nil {
		log.Fatal(err)
	}

	//Track the mixer connection. Functions for onDisconnect and onConnect are passed
	mixer.TrackConnection(
		func() {
			log.Println("Connection considered permanently lost after 20 seconds. Exiting...")
			os.Exit(1)
		},
		func() {
			log.Println("Disconnected")
			qmlRoot.enableBusy()
		},
		func() {
			log.Println("Connected")
			qmlRoot.disableBusy()
		})

	/*	go func() {
			http.ListenAndServe(":8080", nil)
		}()
	*/
	gui.QGuiApplication_Exec()
}
