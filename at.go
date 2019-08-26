package gsm

// TxGeneric is for generic AT commands that take a string possibly with parameters encoded in and expect only an OK or ERROR result code.
type TxGeneric struct {
	AT string
}

// TxCMGS represents a sent message to the network. "Da" is the destination address which is a phone number. "Toda" is the type of destination address which can be 129 for national or 145 for international. "Text" is the body of the text message. Messages are sent with text formating (AT+CMGF=1) so restrictions apply.
type TxCMGS struct {
	Da   string
	Toda int
	Text string
}

// RxCMT represents incoming SMS to the modem which is routed directly to us, not stored on the SIM. We are in text mode (AT+CMGF=1) with detailed header information (AT+CSDH=1). "Oa" is the originator address number which is a phone number. "Length" is the text length in bytes. "Data" is the text of length "Length"
type RxCMT struct {
	Oa     string
	Length int
	Data   string
}

// RxCMTI represents incoming SMS to the modem which is stored on the SIM. "Memr" is the memory storage where the new message is stored. Possible options are SM or ME. "Index" is the location on the memory where SMS is stored.
type RxCMTI struct {
	Memr  string
	Index string
}

// RxCMGS represents a successfully sent message to the network as a result of the command AT+CMGS. "Mr" is the message reference number.
type RxCMGS struct {
	Mr string
}
