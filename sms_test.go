package gsm

import (
	"log"
	"testing"
)

func TestSendSMS(t *testing.T) {
	t.Skip()
	modem, _ := New("serial_tcp")

	go modem.ReadTTY()
	go func() {
		modem.SendSMS("9376721929", "a")
	}()

	for {
		cmti := <-modem.cmti
		log.Printf("%s,%s", cmti.memr, cmti.index)
	}
}
