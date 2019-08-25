package gsm

type TxGeneric struct {
	AT string
}

type TxCMGS struct {
	Da   string
	Toda int
	Text string
}

type RxCMT struct {
	Oa     string
	Scts   string
	Length int
	Data   string
}

type RxCMTI struct {
	Memr  string
	Index string
}

type RxCMGS struct {
	Mr string
}
