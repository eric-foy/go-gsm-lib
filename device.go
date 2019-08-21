package gsm

import (
	"errors"
	"net"

	"github.com/jacobsa/go-serial/serial"
)

// There has to be an already existing interface like this...
type Device interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type Modem struct {
	device  Device
	Results chan string
	Cmt     chan CMT
	Cmti    chan CMTI
	Cmgs    chan CMGS
}

type CMT struct {
	Oa     string
	Scts   string
	Length int
	Data   string
}

type CMTI struct {
	Memr  string
	Index string
}

type CMGS struct {
	Mr string
}

func New(method, device string) (modem *Modem, err error) {
	var dev Device
	switch method {
	// This is serial piped to a tcp server with socat:
	// 	`socat -x FILE:/dev/ttyAMA0 TCP-LISTEN:7875`
	case "serial_tcp":
		dev, err = net.Dial("tcp", device)
		if err != nil {
			return nil, err
		}
	// TTY set with:
	// 	`stty 9600 -F /dev/ttyAMA0 ignpar -icrnl -opost -onlcr -isig -icanon -echo`
	case "serial":
		options := serial.OpenOptions{
			PortName:        device,
			BaudRate:        19200,
			DataBits:        8,
			StopBits:        1,
			MinimumReadSize: 1,
		}
		dev, err = serial.Open(options)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unmatched method")
	}

	modem = &Modem{
		device:  dev,
		Results: make(chan string),
		Cmt:     make(chan CMT),
		Cmti:    make(chan CMTI),
		Cmgs:    make(chan CMGS),
	}
	return modem, nil
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
	// This is used to read correct bytes for incoming SMS.
	modem.AT("AT+CSDH=1")

	// init string
	modem.AT("AT+CNMI=2,2,0,1,0")
}

func (modem *Modem) AT(cmd string) string {
	modem.Write("%s\n", cmd)

	// TODO detecting and wrapping ERROR
	return <-modem.Results
}
