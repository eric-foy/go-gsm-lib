package gsm

import (
	"fmt"
	"strconv"
	"strings"
)

// SendMessage assembles a CMGS message and executes the command. Then writes the text body appending a Ctrl+z to signify the end of the text.
func (modem *Modem) SendMessage(da string, toda int, text string) {
	modem.Write("AT+CMGS=\"%s\",%d\n", da, toda)
	modem.Write("%s%c", text, '\x1A')
}

// ParseCMT parses incoming SMS messages. First locating the byte length of the text body then reading that many bytes.
func (modem *Modem) ParseCMT(line string) RxCMT {
	fields := strings.Split(line[6:], ",")
	length, err := strconv.Atoi(fields[10])
	if err != nil {
		fmt.Println("SMS length conv error")
	}
	data := string(modem.ReadBytes(length))
	return RxCMT{
		Oa:     fields[0][1 : len(fields[0])-1],
		Length: length,
		Data:   data,
	}
}
