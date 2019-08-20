package gsm

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

type CMTI struct {
	memr  string
	index string
}

func (modem *Modem) ReadTTY(results chan []byte, indications chan CMTI) {
	defer modem.connection.Close()

	buffer := make([]byte, 0, 4096)
	temp := make([]byte, 1460) // 1460 is max buffer size for AT
	for {
		bytesRead, err := modem.connection.Read(temp)
		if err != nil {
			fmt.Println("read error:", err)
		}

		buffer = append(buffer, temp[:bytesRead]...)

		// TODO read indication SMS
		newMsgStart := bytes.Index(buffer, []byte("\r\n+CMTI:"))
		if newMsgStart != -1 {
			newMsgEnd := bytes.Index(buffer[newMsgStart+2:], []byte("\r\n"))
			if newMsgEnd != -1 {
				indications <- CMTI{
					memr:  string(buffer[newMsgStart+10 : newMsgStart+12]),
					index: string(buffer[newMsgStart+14 : newMsgEnd+2]),
				}
				buffer = append(buffer[:newMsgStart], buffer[newMsgEnd+4:]...)
			}
		}

		// TODO read suffix ERROR

		if bytes.HasSuffix(buffer, []byte("\r\nOK\r\n")) {
			results <- []byte(buffer)
			buffer = buffer[:0]
		}
	}
}

func (modem *Modem) AT(cmd string, results chan []byte) string {
	modem.SendAT(cmd)
	return modem.ReceiveAT(results)
}

func (modem *Modem) SendAT(cmd string) {
	modem.connection.Write([]byte(cmd + "\n"))
	log.Printf("Sent: %s\n", cmd)
}

func (modem *Modem) ReceiveAT(results chan []byte) string {
	result := string(<-results)
	log.Printf("Received: %s\n", strings.TrimSpace(result))
	return result
}
