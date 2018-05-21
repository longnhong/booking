package common

type ConfigSystem struct {
	LinkCetm           string
	LinkSearchMap      string
	PortBooking        string
	TimeSetCache       int
	CycleDayMissed     string
	KmSearch           float64
	SendNotifyStartDay float64
	SendNotifyBfHour   float64
	StartNear          float64
	EndNear            float64
	StartOut           float64
	EndOut             float64
	ScanNear           float64
	CyclePushDay       int
	CyclePushTicket    int
	UserTicketDay      int
}

var ConfigSystemBooking = ConfigSystem{}
