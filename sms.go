package gsm

import "errors"

func (modem *Modem) SendSMS(number, text string) (CMGS, error) {
	modem.Write("AT+CMGS=\"%s\",145\n", number)
	modem.Write("%s%c", text, '\x1A')

	cmgs := <-modem.Cmgs
	result := <-modem.Results
	if result != "OK" {
		return cmgs, errors.New("problem sending SMS message")
	}
	return cmgs, nil
}
