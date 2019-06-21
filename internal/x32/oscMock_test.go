package x32_test

import (
	"time"

	osc "github.com/vchrisr/go-osc"
)

type oscClientMock struct {
	connected    bool
	sentMessages []*osc.Message
}

func (c *oscClientMock) Connect(retries int) error {
	c.connected = true
	return nil
}

func (c *oscClientMock) Disconnect() error {
	c.connected = false
	return nil
}

func (c *oscClientMock) Send(m *osc.Message) error {
	c.sentMessages = append(c.sentMessages, m)
	return nil
}

func (c *oscClientMock) Receive(t time.Duration) (*osc.Message, error) {
	time.Sleep(100 * time.Millisecond)
	return c.sentMessages[len(c.sentMessages)-1], nil
}
