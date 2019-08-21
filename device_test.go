package gsm

import (
	"fmt"
	"testing"
)

func TestInitDevice(t *testing.T) {
	//t.Skip()
	modem, _ := New("serial_tcp")

	go modem.ReadTTY()
	go modem.InitDevice()
	for {
		cmt := <-modem.cmt
		fmt.Printf("%v", cmt)
	}
}
