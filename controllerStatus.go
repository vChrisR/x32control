package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/therecipe/qt/core"
)

type controllerStatus struct {
	core.QObject

	_ string `property:"ipaddress"`
	_ int    `property:"brightness"`

	_ func()     `constructor:"init"`
	_ func(bool) `signal:"shutdown,auto"`
	_ func(int)  `signal:"changeBrightness,auto"`
}

func (c *controllerStatus) init() {
	if b, err := ioutil.ReadFile("/sys/class/backlight/rpi_backlight/brightness"); err == nil {
		brightness, _ := strconv.Atoi(string(b[:len(b)-1]))
		c.SetBrightness(brightness)
	}

	c.SetIpaddress(c.getMyIp())
}

func (c *controllerStatus) getMyIp() string {
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

func (c *controllerStatus) shutdown(restart bool) {
	if restart {
		fmt.Println("Restarting...")
		exec.Command("sudo", "reboot").Start()
		return
	}

	fmt.Println("Shutting down")
	exec.Command("sudo", "poweroff").Start()
}

func (c *controllerStatus) changeBrightness(brightness int) {
	f, err := os.OpenFile("/sys/class/backlight/rpi_backlight/brightness", os.O_RDWR, os.ModeCharDevice)

	if err != nil {
		return
	}

	_, err = fmt.Fprint(f, strconv.Itoa(brightness))
	if err != nil {
		fmt.Println(err.Error())
	}
}
