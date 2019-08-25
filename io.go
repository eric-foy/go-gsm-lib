package gsm

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (modem *Modem) ReadTTY() {
	buf := bufio.NewReader(modem.Device)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
		}
		line = strings.TrimSpace(line)

		fmt.Println(line)

		switch {
		case line == "OK" || line == "ERROR":
			go func() { modem.RespCode <- line }()
		case len(line) >= 5 && line[:5] == "+CMT:":
			fields := strings.Split(line[6:], ",")
			length, err := strconv.Atoi(fields[10])
			if err != nil {
				fmt.Println("SMS length conv error")
			}
			tmp := make([]byte, length)
			n, err := io.ReadFull(buf, tmp)
			if err != nil {
				fmt.Println("SMS read error:", err)
			}
			fmt.Println(string(tmp[:n]))
			cmt := RxCMT{
				Oa:     fields[0][1 : len(fields[0])-1],
				Scts:   strings.Join(fields[2:4], ","),
				Length: length,
				Data:   string(tmp[:n]),
			}
			go func() { modem.RxAT <- cmt }()
		case len(line) >= 6 && line[:6] == "+CMTI:":
			fields := strings.Split(line[7:], ",")
			cmti := RxCMTI{
				Memr:  fields[0][1:3],
				Index: fields[1],
			}
			go func() { modem.RxAT <- cmti }()
		case len(line) >= 6 && line[:6] == "+CMGS:":
			cmgs := RxCMGS{
				Mr: line[7:],
			}
			go func() { modem.RxAT <- cmgs }()
		}
	}
}

func (modem *Modem) WriteTTY() {
	for {
		switch tx := (<-modem.TxAT).(type) {
		case TxGeneric:
			modem.Write("%s\n", tx.AT)
		case TxCMGS:
			modem.SendMessage(tx.Da, tx.Toda, tx.Text)
		default:
			fmt.Printf("Unmatched AT command: %v\n", tx)
		}

		respCode := <-modem.RespCode
		if respCode != "OK" {
			fmt.Printf("Non ok returned from AT command: %s\n", respCode)
		}
	}
}

func (modem *Modem) Write(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	modem.Device.Write([]byte(fmt.Sprintf(format, a...)))
}
