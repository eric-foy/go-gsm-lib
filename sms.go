package gsm

import (
	"fmt"
	"log"
)

func (modem *Modem) SMS(number, text string, results chan []byte) string {
	modem.connection.Write([]byte(fmt.Sprintf("AT+CMGW=\"%s\"", number)))
	modem.connection.Write([]byte(fmt.Sprintf("%s%c", text, '\x1A')))
	log.Printf("SMS sent:\nTo: %s\nBody: %s\n", number, text)
	return modem.ReceiveAT(results)
}
