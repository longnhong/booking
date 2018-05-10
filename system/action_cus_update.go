package system

import (
	"cetm_booking/o/ticket_onl"
	"encoding/json"
)

func (action *TicketAction) CusUpdate(ticket *ticket_onl.TicketBooking) {
	var data *ticket_onl.TicketUpdate
	err := json.Unmarshal(action.Extra, &data)
	if err != nil {
		action.SetError(err)
		return
	}
	tkUp, err := ticket.UpdateTicketBookingByCustomer(data)
	if err != nil {
		action.SetError(err)
		return
	}
	action.Ticket = tkUp
}
