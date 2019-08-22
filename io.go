package gsm

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (modem *Modem) ReadTTY() {
	defer modem.device.Close()

	buf := bufio.NewReader(modem.device)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
		}
		line = strings.TrimSpace(line)

		fmt.Println(line)

		switch {
		case line == "OK" || line == "ERROR":
			modem.Results <- line
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
			cmt := CMT{
				Oa:     fields[0][1 : len(fields[0])-1],
				Scts:   strings.Join(fields[2:4], ","),
				Length: length,
				Data:   string(tmp[:n]),
			}
			modem.Cmt <- cmt
		case len(line) >= 6 && line[:6] == "+CMTI:":
			fields := strings.Split(line[7:], ",")
			modem.Cmti <- CMTI{
				Memr:  fields[0][1:3],
				Index: fields[1],
			}
		case len(line) >= 6 && line[:6] == "+CMGS:":
			modem.Cmgs <- CMGS{
				Mr: line[7:],
			}
		}
	}
}

func (modem *Modem) Write(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	modem.device.Write([]byte(fmt.Sprintf(format, a...)))
}
