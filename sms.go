package gsm

func (modem *Modem) SendSMS(number, text string) CMGS {
	modem.Write("AT+CMGS=\"%s\",129\n", number)
	modem.Write("%s%c", text, '\x1A')
	return <-modem.cmgs
}

func (modem *Modem) ReadMessage(index string) string {
	modem.Write("AT+CMGR=%s\n", index)
	return ""
}
