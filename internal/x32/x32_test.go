package x32_test

import (
	"testing"
	"time"

	"github.com/vchrisr/x32control/internal/x32"
)

//
func TestRecallScene(t *testing.T) {
	c := &oscClientMock{}
	mixer := x32.New(c)
	err := mixer.RecallScene(2)
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}
	time.Sleep(600 * time.Millisecond) //function under test wait 500 ms to give board time to recall

	oscAddresses := []string{"/-action/goscene", "/xremote"}

	for i, addrs := range oscAddresses {
		if c.sentMessages[i].Address != addrs {
			t.Errorf("Want osc Address %v, have osc Address: %v", addrs, c.sentMessages[i].Address)
		}
	}

	if c.sentMessages[0].Arguments[0].(int32) != 2 {
		t.Errorf("Want osc Argument %v, have osc Argument: %v", 2, c.sentMessages[0].Arguments)
	}
}

func TestStart(t *testing.T) {
	c := &oscClientMock{}
	mixer := x32.New(c)
	err := mixer.Start()
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}

	if !c.connected {
		t.Errorf("Start is supposed to connect but it didn't")
	}

	expected := []string{"/info", "/xremote", "/batchsubscribe"}

	for i, addr := range expected {
		if c.sentMessages[0].Address != "/info" {
			t.Errorf("message index %v expected: %v but got: %v", i, addr, c.sentMessages[i].Address)
		}
	}
}
