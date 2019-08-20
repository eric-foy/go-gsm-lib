package gsm

import (
	"testing"
)

func TestInitDevice(t *testing.T) {
	m, _ := New("serial_tcp")
	m.InitDevice()
}
