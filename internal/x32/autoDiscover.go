package x32

import (
	"fmt"
	"time"

	osc "github.com/vchrisr/go-osc"
)

func AutoDiscover(retries int) ([]string, error) {
	discoveries := make(map[string]*osc.Message)
	var err error

	//retries needed if run directly after device bootup
	for retries > 0 {
		discoveries, err = osc.AutoDiscover(10023, osc.NewMessage("/info"))
		if err == nil && len(discoveries) > 0 {
			break
		}

		time.Sleep(1 * time.Second)
		retries--
	}

	if len(discoveries) == 0 {
		err = fmt.Errorf("No x32 discovered on the network. Error: %v", err)
	}

	if err != nil {
		return nil, err
	}

	var validDiscoveries []string
	for ip, msg := range discoveries {
		if msg.Address == "/info" && msg.Arguments[2].(string) == "X32" {
			validDiscoveries = append(validDiscoveries, ip)
		}
	}

	if len(validDiscoveries) == 0 {
		return nil, fmt.Errorf("No X32 discovered")
	}

	return validDiscoveries, nil

}
