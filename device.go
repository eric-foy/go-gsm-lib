package gsm

import (
	"bufio"
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
	Device   Device
	Reader   *bufio.Reader
	RespCode chan string
	RxAT     chan interface{}
	TxAT     chan interface{}
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
			BaudRate:        9600,
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
		Device:   dev,
		Reader:   bufio.NewReader(dev),
		RespCode: make(chan string),
		TxAT:     make(chan interface{}),
		RxAT:     make(chan interface{}),
	}
	return modem, nil
}

func (modem *Modem) InitDevice() {
	// reset the modem
	modem.TxAT <- TxGeneric{AT: "ATZ"}

	// reset settings
	modem.TxAT <- TxGeneric{AT: "AT&F"}

	// echo off
	modem.TxAT <- TxGeneric{AT: "ATE0"}

	// text mode
	modem.TxAT <- TxGeneric{AT: "AT+CMGF=1"}

	// detailed header information in text mode
	// This is used to read correct bytes for incoming SMS.
	modem.TxAT <- TxGeneric{AT: "AT+CSDH=1"}

	// init string
	modem.TxAT <- TxGeneric{AT: "AT+CNMI=2,2,0,1,0"}
}
