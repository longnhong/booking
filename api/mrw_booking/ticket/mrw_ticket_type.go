package ticket

import (
	"cetm_booking/o/auth/user"
	"cetm_booking/o/ticket_onl"
)

type DataBank struct {
	Data   *InfoBank `json:"data"`
	Status string    `json:"status"`
}

type InfoBank struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Address     string  `json:"address"`
	BranchID    string  `json:branch_id`
	CountPeople int     `json:"count_people"`
}

type DataTicketCetm struct {
	Id   string `json:"id"`
	Cnum string `json:"cnum"`
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
