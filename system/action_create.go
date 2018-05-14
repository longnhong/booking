package system

import (
	"cetm_booking/o/ticket_onl"
	"encoding/json"
	"fmt"
)

func (action *TicketAction) actionCreate() {
	var data = struct {
		PushToken string                         `json:"push_token"`
		Ticket    ticket_onl.TicketBookingCreate `json:"ticket"`
	}{}
	err := json.Unmarshal(action.Extra, &data)
	fmt.Printf("Unmarshal", data)
	if err != nil {
		action.SetError(err)
		return
	}
	ticket, err := data.Ticket.CrateTicketBooking()
	fmt.Printf("CrateTicketBooking", err)
	if err != nil {
		action.SetError(err)
		return
	}
	action.Ticket = ticket
	sendFee(data.PushToken, ticket)
}
