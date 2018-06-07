package system

import (
	"cetm_booking/o/push_token"
	"cetm_booking/o/ticket_onl"
	"encoding/json"
	"errors"
)

func (action *TicketAction) actionFinish(ticket *ticket_onl.TicketBooking) {
	var data *ticket_onl.UpdateCetm
	var err = json.Unmarshal(action.Extra, &data)
	if err != nil {
		action.SetError(err)
		return
	}
	var status = ticket.Status
	if status == ticket_onl.BookingStateSancelled || ticket_onl.BookingStateFinished == status {
		action.SetError(errors.New("Vé đã được phản hồi"))
		return
	}
	err = ticket.UpdateTicketBookingByCetm(data)
	if err != nil {
		action.SetError(err)
		return
	}
	action.Ticket = ticket

	if action.Action == ticket_onl.BookingStateFinished {
		pDevices, err := push_token.GetPushsUserId(ticket.CustomerID)
		if err != nil {
			action.SetError(err)
			return
		}
		sendFeedback(pDevices, ticket, action.Action)
	}
}
