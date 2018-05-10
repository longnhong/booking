package ctrl_to_cetm

import (
	"cetm_booking/o/auth/user"
	"cetm_booking/o/ticket_onl"
)

type InfoBank struct {
	Lat                 float64   `json:"lat"`
	Lng                 float64   `json:"lng"`
	Address             string    `json:"address"`
	KmScan              float64   `json:"km_scan"`
	BranchID            string    `json:"branch_id"`
	NameBranch          string    `json:"name"`
	CountPeople         int       `json:"count_people"`
	CountCounterService int       `json:"count_counter_service"`
	MaxServingMinute    int32     `json:"max_serving_minute"`
	ServiceInCounters   []string  `json:"service_in_counters"`
	Counters            []Counter `json:"counters"`
}

type Counter struct {
	BranchID  string   `json:"branch_id"`
	Cnum      string   `json:"cnum"`
	Code      string   `json:"code"`
	DevAddr   int      `json:"dev_addr"`
	Dtime     int      `json:"dtime"`
	ID        string   `json:"id"`
	Mtime     int      `json:"mtime"`
	Name      string   `json:"name"`
	Pservices []string `json:"pservices"`
	Services  []string `json:"services"`
	Vservices []string `json:"vservices"`
}

type DataTicketCetm struct {
	Id          string `json:"id"`
	Cnum        string `json:"cnum"`
	CountPeople int    `json:"count_people"`
}

type DataTicketBookNow struct {
	Data          DataTicketCetm           `json:"data"`
	TicketBooking ticket_onl.TicketBooking `json:"ticket_booking"`
	Status        string                   `json:"status"`
}

type DataTicketSendCetm struct {
	TicketBooking *ticket_onl.TicketBooking `json:"ticket_booking"`
	Customer      *user.User                `json:"customer"`
}
