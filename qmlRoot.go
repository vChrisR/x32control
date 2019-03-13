package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
)

type QmlRoot struct {
	core.QObject

	chStrips ChannelStrips
	mixer    *x32

	_ bool   `property:"busy"`
	_ string `property:"ipaddress"`
	_ int    `property:"brightness"`

	_ func()                      `constructor:"init"`
	_ func(bool)                  `slot:"shutdown,auto"`
	_ func(int)                   `slot:"changeBrightness,auto"`
	_ func(int)                   `slot:"recallClicked,auto"`
	_ func(string, *core.QObject) `slot:"registerChannelStrip,auto"`
	_ func(string, int)           `signal:"receiveFaderValue"`
	_ func(string, float32)       `slot:"sendFaderValue,auto"`
	_ func(string, bool)          `slot:"sendMute,auto"`
}

func initQmlRoot(view *quick.QQuickView, conf config, mixer *x32) *QmlRoot {
	q := NewQmlRoot(nil)
	q.mixer = mixer
	q.chStrips = make(ChannelStrips)

	confJson, _ := json.Marshal(conf)

	view.RootContext().SetContextProperty2("controllerConfig", core.NewQVariant14(string(confJson)))
	view.RootContext().SetContextProperty("QmlRoot", q)

	return q
}

func (q *QmlRoot) init() {
	q.SetBusy(true)

	if b, err := ioutil.ReadFile("/sys/class/backlight/rpi_backlight/brightness"); err == nil {
		brightness, _ := strconv.Atoi(string(b[:len(b)-1]))
		q.SetBrightness(brightness)
	}

	//Get and set the IP address every 5 seconds.
	//This will update the IP properly even if the app was started when network was not yet available.
	//Also account for changed ip address (DHCP)
	go func() {
		for {
			q.SetIpaddress(q.getMyIp())
			time.Sleep(5 * time.Second)
		}
	}()
}

func (q *QmlRoot) getMyIp() string {
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if !ip.IsLoopback() && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return ""

}

func (q *QmlRoot) shutdown(restart bool) {
	if restart {
		log.Println("Restarting...")
		exec.Command("sudo", "reboot").Start()
		return
	}

	log.Println("Shutting down")
	exec.Command("sudo", "poweroff").Start()
}

func (q *QmlRoot) recallClicked(scene int) {
	if err := q.mixer.RecallScene(int(scene)); err != nil {
		log.Println(err.Error())
	}

	time.Sleep(500 * time.Millisecond)

	for _, channel := range q.chStrips {
		channel.updateFromMixer()
	}
}

func (q *QmlRoot) enableBusy() {
	q.SetBusy(true)
	for _, chStrip := range q.chStrips {
		chStrip.lastFaderPosition = 0
	}
}

func (q *QmlRoot) disableBusy() {
	q.SetBusy(false)
	for _, chStrip := range q.chStrips {
		chStrip.updateFromMixer()
	}

	q.mixer.Send(osc.NewMessage("/xremote"))
	q.mixer.RequestMetering()
}

func (q *QmlRoot) changeBrightness(brightness int) {
	f, err := os.OpenFile("/sys/class/backlight/rpi_backlight/brightness", os.O_RDWR, os.ModeCharDevice)

	if err != nil {
		return
	}

	_, err = fmt.Fprint(f, strconv.Itoa(brightness))
	if err != nil {
		log.Println(err.Error())
	}
}

func (q *QmlRoot) registerChannelStrip(addr string, qmlObj *core.QObject) {
	q.chStrips[addr] = &ChannelStrip{
		qmlObj:            qmlObj,
		mixerChannel:      NewX32Channel(addr, q.mixer),
		lastFaderPosition: 0,
	}
}

func (q *QmlRoot) sendFaderValue(address string, pos float32) {
	q.chStrips[address].sendFaderPosition(pos)
}

func (q *QmlRoot) sendMute(address string, checked bool) {
	q.chStrips[address].mixerChannel.setMute(checked)
}
