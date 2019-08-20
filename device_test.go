package gsm

import (
	"log"
	"testing"
)

func TestInitDevice(t *testing.T) {
	modem, _ := New("serial_tcp")

	results := make(chan []byte)
	indications := make(chan CMTI)
	go modem.ReadTTY(results, indications)

	modem.InitDevice(results)
	for {
		cmti := <-indications
		log.Printf("%s,%s", cmti.memr, cmti.index)
	}
}
