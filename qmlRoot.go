package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
)

type QmlRoot struct {
	core.QObject

	_ bool   `property:"busy"`
	_ string `property:"ipaddress"`
	_ int    `property:"brightness"`

	_ func() `constructor:"init"`

	_ func(bool) `slot:"shutdown,auto"`
	_ func(int)  `slot:"changeBrightness,auto"`
	_ func(int)  `slot:"recallClicked"`
}

func initQmlRoot(view *quick.QQuickView, conf config) *QmlRoot {
	q := NewQmlRoot(nil)

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
		fmt.Println("Restarting...")
		exec.Command("sudo", "reboot").Start()
		return
	}

	fmt.Println("Shutting down")
	exec.Command("sudo", "poweroff").Start()
}

func (q *QmlRoot) changeBrightness(brightness int) {
	f, err := os.OpenFile("/sys/class/backlight/rpi_backlight/brightness", os.O_RDWR, os.ModeCharDevice)

	if err != nil {
		return
	}

	_, err = fmt.Fprint(f, strconv.Itoa(brightness))
	if err != nil {
		fmt.Println(err.Error())
	}
}
