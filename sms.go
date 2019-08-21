package gsm

func (modem *Modem) SendSMS(number, text string) {
	modem.Write("AT+CMGS=\"%s\",129\n", number)
	modem.Write("%s%c", text, '\x1A')
}

func (modem *Modem) WriteMsgToMem(number, text string) {
}

func (modem *Modem) SendMsgFromStorage(index string) {
}

func (modem *Modem) ReadMessage(index string) string {
	modem.Write("AT+CMGR=%s\n", index)
	return ""
}
