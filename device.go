package gsm

import (
	"bufio"
	"errors"
	"net"

	serial "github.com/jacobsa/go-serial/serial"
)

// Device enables the New method to take a generic device that can Read Write and Close. Serial devices and TCP connections both satisfy this.
type Device interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

// Modem represents communication with the GSM modem. It includes "Device" which allows us to defer closing until ready and access the underlying Read and Write. "Reader" which gives us buffered reading on the special device. "RespCode" which is a channel that communicates success or failure of an executed AT command. "RxAT" for received and parsed AT commands. "TxAT" that takes TxAT structures, formats them, and sends them to the modem.
type Modem struct {
	Device   Device
	Reader   *bufio.Reader
	RespCode chan string
	RxAT     chan interface{}
	TxAT     chan interface{}
}

// New creates a modem and makes the channels it uses. It can be a serial piped to a tcp server with socat. `socat -x FILE:/dev/ttyAMA0 TCP-LISTEN:7875` (serial_tcp). It can also be a serial device such as /dev/ttyAMA0. I set my device with flags as given to stty. `stty 9600 -F /dev/ttyAMA0 ignpar -icrnl -opost -onlcr -isig -icanon -echo`.
func New(method, device string) (modem *Modem, err error) {
	var dev Device
	switch method {
	case "serial_tcp":
		dev, err = net.Dial("tcp", device)
		if err != nil {
			return nil, err
		}
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

// InitDevice executes commands to set up the device so it behaves as expected. "ATZ" resets the modem. "AT&F" resets settings. "ATE0" turns echoing of commands off. "AT+CMGF=1" turns text mode on. "AT+CSDH=1" sets detailed header information in text mode, this is used to read correct bytes for incoming SMS. "AT+CNMI=2,2,0,1,0" sets up how received indications are communicated such as if they are stored on the SIM or sent directly to the serial. AT+CNMI should be looked up for farther understanding, its a long one.
func (modem *Modem) InitDevice() {
	modem.TxAT <- TxGeneric{AT: "ATZ"}
	modem.TxAT <- TxGeneric{AT: "AT&F"}
	modem.TxAT <- TxGeneric{AT: "ATE0"}
	modem.TxAT <- TxGeneric{AT: "AT+CMGF=1"}
	modem.TxAT <- TxGeneric{AT: "AT+CSDH=1"}
	modem.TxAT <- TxGeneric{AT: "AT+CNMI=2,2,0,1,0"}
}
