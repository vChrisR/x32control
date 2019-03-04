package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

const secondsFrom1900To1970 = 2208988800

type oscHandler func(*osc.Message)

type x32 struct {
	ip          string
	port        int
	connection  *net.UDPConn
	connected   bool
	lastMessage time.Time
	oscHandlers map[string]oscHandler
}

func (x *x32) connect() error {
	fmt.Printf("Connecting to %v...\n", x.ip)
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", x.ip, x.port))
	if err != nil {
		return err
	}

	retries := 8
	for x.connection == nil && retries > 0 {
		if x.connection, err = net.DialUDP("udp", nil, addr); err != nil {
			fmt.Println(err.Error())
			time.Sleep(2 * time.Second)
		}
		retries--
	}

	if x.connection == nil {
		return fmt.Errorf("Unable to connect to the network.")
	}

	return nil
}

func (x *x32) Start() error {
	if err := x.connect(); err != nil {
		return err
	}

	if err := x.GetInfo(); err != nil {
		return err
	}
	msgStream := x.listen()
	x.startDispatcher(msgStream)
	x.StartXremote()
	x.StartMetering()

	return nil
}

func (x *x32) Disconnect() {
	x.UnSubscribe()
	x.connection.Close()
}

func (x *x32) Send(msg *osc.Message) error {
	if x.connection == nil {
		return fmt.Errorf("Unable to send, not connected.")
	}

	data, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	if _, err = x.connection.Write(data); err != nil {
		return err
	}

	return nil
}

func (x *x32) Receive(timeout time.Duration) (error, *osc.Message) {
	buffer := make([]byte, 65535)
	x.connection.SetReadDeadline(time.Now().Add(timeout))
	_, _, err := x.connection.ReadFrom(buffer)
	if err != nil {
		return err, nil
	}

	var start int
	msg, _, err := readMessage(bufio.NewReader(bytes.NewBuffer(buffer)), &start)
	if err != nil {
		return err, nil
	}

	//fmt.Println(msg.String())

	return nil, msg
}

func (x *x32) RecallScene(sceneNr int) error {
	msgAddr := "/-action/goscene"
	msg := osc.NewMessage(msgAddr, int32(sceneNr))
	return x.Send(msg)
}

func (x *x32) GetInfo() error {
	msg := osc.NewMessage("/info")
	return x.Send(msg)
}

func (x *x32) StartXremote() {
	xremoteMsg := osc.NewMessage("/xremote")
	x.Send(xremoteMsg)

	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()

		for _ = range tick.C {
			if x.connected {
				x.Send(xremoteMsg)
			}
		}
	}()
}

func (x *x32) RequestMetering() {
	msg := osc.NewMessage("/batchsubscribe")
	msg.Append("/metering")
	msg.Append("/meters/0")
	msg.Append(int32(0))
	msg.Append(int32(0))
	msg.Append(int32(0))

	x.Send(msg)

}

func (x *x32) StartMetering() {
	x.RequestMetering()

	go func() {
		tick := time.NewTicker(9 * time.Second)
		defer tick.Stop()
		renewMsg := osc.NewMessage("/renew")
		renewMsg.Append("/metering")
		for _ = range tick.C {
			if x.connected {
				x.Send(renewMsg)
			}
		}
	}()
}

func (x *x32) UnSubscribe() {
	x.Send(osc.NewMessage("/unsubscribe"))
}

func (x *x32) listen() <-chan *osc.Message {
	stream := make(chan *osc.Message, 16)

	go func() {
		defer close(stream)

		for {
			err, msg := x.Receive(500 * time.Millisecond)
			if err != nil {
				continue
			}

			x.lastMessage = time.Now()

			select {
			case stream <- msg:
			default:
				fmt.Println("Receive buffer overrun, message discarded")
			}
		}
	}()
	/*
		go func() {
			data := MeterData{
				NumberOfFloats: 70,
				Channel:        [32]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
				Aux:            [8]float32{0.75, 0.75, 0.75, 0.75, 0.75, 0.75, 0.75, 0.75},
				FxReturn: [4]struct {
					Left  float32
					Right float32
				}{
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5},
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5},
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5},
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5}},
				Bus:    [16]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
				Matrix: [6]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			}

			data2 := MeterData{
				NumberOfFloats: 70,
				Channel:        [32]float32{0.75, 0.75, 0.75, 0.75, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.75, 0.75, 0.75, 0.75, 0.75, 0.75},
				Aux:            [8]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
				FxReturn: [4]struct {
					Left  float32
					Right float32
				}{
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5},
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5},
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5},
					struct {
						Left  float32
						Right float32
					}{Left: 0.5, Right: 0.5}},
				Bus:    [16]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
				Matrix: [6]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			}

			var buf, buf2 bytes.Buffer
			binary.Write(&buf, binary.LittleEndian, data)
			binary.Write(&buf2, binary.LittleEndian, data2)

			msg1 := osc.NewMessage("/metering")
			msg1.Append(buf.Bytes())
			msg2 := osc.NewMessage("/metering")
			msg2.Append(buf2.Bytes())
			for {
				stream <- msg1
				time.Sleep(50 * time.Millisecond)
				stream <- msg2
				time.Sleep(50 * time.Millisecond)
			}
		}()*/

	return stream
}

func (x *x32) TrackConnection(connLost func(), connected func()) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			for _ = range ticker.C {
				if time.Since(x.lastMessage).Seconds() > 1 {
					err := x.GetInfo()
					if err != nil {
						fmt.Println(err.Error())
					}

					if time.Since(x.lastMessage).Seconds() > 3 && x.connected {
						x.connected = false
						go connLost()
					}
				} else {
					if !x.connected {
						x.connected = true
						go connected()
					}
				}
			}
		}
	}()
}

func (x *x32) Handle(route string, handler oscHandler) {
	x.oscHandlers[route] = handler
}

func (x *x32) startDispatcher(msgChan <-chan *osc.Message) {
	go func() {
		for msg := range msgChan {
			addrElements := strings.Split(msg.Address, "/")
			if handler, exists := x.oscHandlers[addrElements[1]]; exists {
				handler(msg)
			}
		}
	}()
}

func NewX32(ip string, port int) *x32 {
	return &x32{
		ip:          ip,
		port:        port,
		connected:   false,
		oscHandlers: make(map[string]oscHandler),
	}
}

type x32channel struct {
	mixer    *x32
	baseAddr string
}

func (x x32channel) getMsgAddress(action string) string {
	var msgAddr string
	if strings.Split(x.baseAddr, "/")[1] == "dca" {
		msgAddr = fmt.Sprintf("%v/%v", x.baseAddr, action)
	} else {
		msgAddr = fmt.Sprintf("%v/mix/%v", x.baseAddr, action)
	}

	return msgAddr
}

func (x x32channel) setMute(muted bool) error {
	var value int32
	value = 0
	if !muted {
		value = 1
	}

	msg := osc.NewMessage(x.getMsgAddress("on"), value)
	return x.mixer.Send(msg)
}

func (x x32channel) getMute() error {
	msg := osc.NewMessage(x.getMsgAddress("on"))
	return x.mixer.Send(msg)
}

func (x x32channel) getName() error {
	msgAddr := fmt.Sprintf("%v/config/name", x.baseAddr)

	msg := osc.NewMessage(msgAddr)
	return x.mixer.Send(msg)
}

func (x x32channel) setFaderPosition(position float32) error {
	msg := osc.NewMessage(x.getMsgAddress("fader"), position)
	return x.mixer.Send(msg)
}

func (x x32channel) getFaderPosition() error {
	msg := osc.NewMessage(x.getMsgAddress("fader"))
	return x.mixer.Send(msg)
}

func NewX32Channel(oscAddress string, mixer *x32) x32channel {
	channel := x32channel{
		mixer:    mixer,
		baseAddr: oscAddress,
	}

	return channel
}

// readMessage from `reader`.
func readMessage(reader *bufio.Reader, start *int) (*osc.Message, bool, error) {
	buf, err := reader.Peek(1)
	if err != nil {
		return nil, false, err
	}

	if buf[0] == '/' {
		// First, read the OSC address
		addr, n, err := readPaddedString(reader)
		if err != nil {
			return nil, true, err
		}
		*start += n

		// Read all arguments
		msg := osc.NewMessage(addr)
		if err = readArguments(msg, reader, start); err != nil {
			return nil, true, err
		}
		return msg, true, nil
	}
	return &osc.Message{}, false, fmt.Errorf("Not an OSC message")
}

// readPaddedString reads a padded string from the given reader. The padding
// bytes are removed from the reader.
func readPaddedString(reader *bufio.Reader) (string, int, error) {
	// Read the string from the reader
	str, err := reader.ReadString(0)
	if err != nil {
		return "", 0, err
	}
	n := len(str)

	// Remove the string delimiter, in order to calculate the right amount
	// of padding bytes
	str = str[:len(str)-1]

	// Remove the padding bytes
	padLen := padBytesNeeded(len(str)) - 1
	if padLen > 0 {
		n += padLen
		padBytes := make([]byte, padLen)
		if _, err = reader.Read(padBytes); err != nil {
			return "", 0, err
		}
	}

	return str, n, nil
}

// readArguments from `reader` and add them to the OSC message `msg`.
func readArguments(msg *osc.Message, reader *bufio.Reader, start *int) error {
	// Read the type tag string
	var n int
	typetags, n, err := readPaddedString(reader)
	if err != nil {
		return err
	}
	*start += n

	// If the typetag doesn't start with ',', it's not valid
	if typetags[0] != ',' {
		return errors.New("unsupported type tag string")
	}

	// Remove ',' from the type tag
	typetags = typetags[1:]

	for _, c := range typetags {
		switch c {
		default:
			return fmt.Errorf("unsupported type tag: %c", c)

		case 'i': // int32
			var i int32
			if err = binary.Read(reader, binary.BigEndian, &i); err != nil {
				return err
			}
			*start += 4
			msg.Append(i)

		case 'h': // int64
			var i int64
			if err = binary.Read(reader, binary.BigEndian, &i); err != nil {
				return err
			}
			*start += 8
			msg.Append(i)

		case 'f': // float32
			var f float32
			if err = binary.Read(reader, binary.BigEndian, &f); err != nil {
				return err
			}
			*start += 4
			msg.Append(f)

		case 'd': // float64/double
			var d float64
			if err = binary.Read(reader, binary.BigEndian, &d); err != nil {
				return err
			}
			*start += 8
			msg.Append(d)

		case 's': // string
			// TODO: fix reading string value
			var s string
			if s, _, err = readPaddedString(reader); err != nil {
				return err
			}
			*start += len(s) + padBytesNeeded(len(s))
			msg.Append(s)

		case 'b': // blob
			var buf []byte
			var n int
			if buf, n, err = readBlob(reader); err != nil {
				return err
			}
			*start += n
			msg.Append(buf)

		case 't': // OSC time tag
			var tt uint64
			if err = binary.Read(reader, binary.BigEndian, &tt); err != nil {
				return nil
			}
			*start += 8
			msg.Append(osc.NewTimetagFromTimetag(tt))

		case 'T': // true
			msg.Append(true)

		case 'F': // false
			msg.Append(false)
		}
	}

	return nil
}

// padBytesNeeded determines how many bytes are needed to fill up to the next 4
// byte length.
func padBytesNeeded(elementLen int) int {
	return 4*(elementLen/4+1) - elementLen
}

// readBlob reads an OSC blob from the blob byte array. Padding bytes are
// removed from the reader and not returned.
func readBlob(reader *bufio.Reader) ([]byte, int, error) {
	// First, get the length
	var blobLen int32
	if err := binary.Read(reader, binary.BigEndian, &blobLen); err != nil {
		return nil, 0, err
	}
	n := 4 + int(blobLen)

	// Read the data
	blob := make([]byte, blobLen)
	if _, err := reader.Read(blob); err != nil {
		return nil, 0, err
	}

	// Remove the padding bytes
	numPadBytes := padBytesNeeded(int(blobLen))
	if numPadBytes > 0 {
		n += numPadBytes
		dummy := make([]byte, numPadBytes)
		if _, err := reader.Read(dummy); err != nil {
			return nil, 0, err
		}
	}

	return blob, n, nil
}
