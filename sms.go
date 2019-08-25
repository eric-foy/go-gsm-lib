package gsm

func (modem *Modem) SendMessage(da string, toda int, text string) {
	modem.Write("AT+CMGS=\"%s\",%d\n", da, toda)
	modem.Write("%s%c", text, '\x1A')
}
