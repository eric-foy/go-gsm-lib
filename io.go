package gsm

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type CMTI struct {
	memr  string
	index string
}

func (modem *Modem) ReadTTY() {
	defer modem.connection.Close()

	buf := bufio.NewReader(modem.connection)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
		}
		line = strings.TrimSpace(line)

		switch {
		case line == "OK" || line == "ERROR":
			modem.results <- line
		case len(line) >= 6 && line[:6] == "+CMTI:":
			fields := strings.Split(line[7:], ",")
			modem.indications <- CMTI{
				memr:  fields[0][1:3],
				index: fields[1],
			}
		}
	}
}

func (modem *Modem) Write(format string, a ...interface{}) {
	log.Printf(format, a...)
	modem.connection.Write([]byte(fmt.Sprintf(format, a...)))
}
