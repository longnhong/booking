package system

import (
	"cetm_booking/o/ticket_onl"
)

func (action *TicketAction) actionDelete(ticket *ticket_onl.TicketBooking) {
	err := ticket.MarkDeleteTicket()
	if err != nil {
		action.SetError(err)
		return
	}
	ticket.Status = ticket_onl.BookingStateDelete
	action.Ticket = ticket
}
