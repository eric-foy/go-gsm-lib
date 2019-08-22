package gsm

import (
	"testing"
	"time"
)

func TestSendSMS(t *testing.T) {
	//t.Skip()
	modem, _ := New("serial_tcp", "192.168.1.130:7875")

	go modem.ReadTTY()
	modem.InitDevice()
	for {
		modem.SendSMS("+19374745540", "a")
		time.Sleep(time.Second * 2)
	}
}
