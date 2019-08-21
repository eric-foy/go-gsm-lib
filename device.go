package gsm

import (
	"errors"
	"net"
)

type Modem struct {
	connection net.Conn
	queue      chan string
	results    chan string
	cmt        chan CMT
	cmti       chan CMTI
	cmgs       chan CMGS
}

type CMT struct {
	oa     string
	scts   string
	length int
	data   string
}

type CMTI struct {
	memr  string
	index string
}

type CMGS struct {
	mr string
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
			queue:      make(chan string),
			results:    make(chan string),
			cmt:        make(chan CMT),
			cmti:       make(chan CMTI),
			cmgs:       make(chan CMGS),
		}
		return modem, nil
	// TTY set with:
	// 	`stty 9600 -F /dev/ttyAMA0 ignpar -icrnl -opost -onlcr -isig -icanon -echo`
	case "serial":
		// TODO connect with /dev/ttyAMA0
	}

	return nil, errors.New("unmatched device")
}

func (modem *Modem) InitDevice() {
	// reset the modem
	modem.AT("ATZ")

	// reset settings
	modem.AT("AT&F")

	// echo off
	modem.AT("ATE0")

	// text mode
	modem.AT("AT+CMGF=1")

	// detailed header information in text mode
	modem.AT("AT+CSDH=1")

	// init string
	modem.AT("AT+CNMI=2,2,0,1,0")
}

func (modem *Modem) AT(cmd string) string {
	modem.Write("%s\n", cmd)

	// TODO detecting and wrapping ERROR
	return <-modem.results
}
