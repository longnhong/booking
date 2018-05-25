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
	if ticket.Status == ticket_onl.BOOKING_STATE_CANCELLED || ticket_onl.BOOKING_STATE_FINISHED == ticket.Status {
		action.SetError(errors.New("Vé đã được phản hồi"))
		return
	}
	err = ticket.UpdateTicketBookingByCetm(data)
	if err != nil {
		action.SetError(err)
		return
	}
	action.Ticket = ticket

	pDevices, err := push_token.GetPushsUserId(ticket.CustomerID)
	if err != nil {
		action.SetError(err)
		return
	}
	sendFeedback(pDevices, ticket, action.Action)
}
