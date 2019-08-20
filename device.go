package gsm

import (
	"errors"
	"net"
)

type Modem struct {
	connection net.Conn
}

func New(device string) (*Modem, error) {
	switch device {
	// This is serial piped to a tcp server with socat:
	// 	`socat -x FILE:/dev/ttyAMA0 TCP-LISTEN:7875`
	case "serial_tcp":
		connection, err := net.Dial("tcp", "192.168.1.130:7875")
		if err != nil {
			return nil, err
		}

		modem := &Modem{
			connection: connection,
		}
		return modem, nil
	// TTY set with:
	// 	`stty 9600 -F /dev/ttyAMA0 ignpar -icrnl -opost -onlcr -isig -icanon -echo`
	case "serial":
		// TODO connect with /dev/ttyAMA0
	}

	return nil, errors.New("unmatched device")
}

func (modem *Modem) InitDevice(results chan []byte) {
	// reset the modem
	modem.AT("ATZ", results)

	// reset settings
	modem.AT("AT&F", results)

	// echo off
	modem.AT("ATE0", results)

	// text mode
	modem.AT("AT+CMGF=1", results)

	// Init string
	modem.AT("AT+CNMI=2,1,0,1,0", results)
}
