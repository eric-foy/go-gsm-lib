package gsm

import (
	"fmt"
	"testing"
)

func TestInitDeviceTCP(t *testing.T) {
	t.Skip()
	modem, _ := New("serial_tcp", "192.168.1.130:7875")

	go modem.ReadTTY()
	go modem.InitDevice()
	for {
		cmt := <-modem.Cmt
		fmt.Printf("%v", cmt)
	}
}

func TestInitDeviceSerial(t *testing.T) {
	t.Skip()
	modem, _ := New("serial", "/dev/ttyAMA0")

	go modem.ReadTTY()
	go modem.InitDevice()
	for {
		cmt := <-modem.Cmt
		fmt.Printf("%v", cmt)
	}
}
