package gsm

import (
	"fmt"
	"strconv"
	"strings"
)

func (modem *Modem) SendMessage(da string, toda int, text string) {
	modem.Write("AT+CMGS=\"%s\",%d\n", da, toda)
	modem.Write("%s%c", text, '\x1A')
}

func (modem *Modem) ParseCMT(line string) RxCMT {
	fields := strings.Split(line[6:], ",")
	length, err := strconv.Atoi(fields[10])
	if err != nil {
		fmt.Println("SMS length conv error")
	}
	data := string(modem.ReadBytes(length))
	return RxCMT{
		Oa:     fields[0][1 : len(fields[0])-1],
		Scts:   strings.Join(fields[2:4], ","),
		Length: length,
		Data:   data,
	}
}
