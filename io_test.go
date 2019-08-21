package gsm

import (
	"log"
	"testing"
	"time"
)

func TestSpamAT(t *testing.T) {
	t.Skip()
	modem, _ := New("serial_tcp")

	go modem.ReadTTY()
	go func() {
		for {
			modem.AT("AT")
			time.Sleep(time.Second)
		}
	}()

	for {
		cmti := <-modem.cmti
		log.Printf("%s,%s", cmti.memr, cmti.index)
	}
}
