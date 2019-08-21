package gsm

func (modem *Modem) SendSMS(number, text string) CMGS {
	modem.Write("AT+CMGS=\"%s\",129\n", number)
	modem.Write("%s%c", text, '\x1A')
	return <-modem.cmgs
}
