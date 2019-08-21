package gsm

import (
	"log"
	"testing"
)

func TestInitDevice(t *testing.T) {
	t.Skip()
	modem, _ := New("serial_tcp")

	go modem.ReadTTY()
	go modem.InitDevice()
	for {
		cmti := <-modem.indications
		log.Printf("%s,%s", cmti.memr, cmti.index)
	}
}
