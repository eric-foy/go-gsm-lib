package gsm

import (
	"log"
	"testing"
	"time"
)

func TestSpamAT(t *testing.T) {
	t.Skip()
	modem, _ := New("serial_tcp", "192.168.1.130:7875")

	go modem.ReadTTY()
	go func() {
		for {
			modem.AT("AT")
			time.Sleep(time.Second)
		}
	}()

	for {
		cmti := <-modem.Cmti
		log.Printf("%s,%s", cmti.Memr, cmti.Index)
	}
}
