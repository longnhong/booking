package system

import (
	"cetm_booking/o/ticket_onl"
	"errors"
)

func (action *TicketAction) handlerAction(ticket *ticket_onl.TicketBooking) {
	switch action.Action {
	case ticket_onl.BookingStateCreated:
		action.actionCreate()
	//case ticket_onl.BookingStateConfirmed:
	case ticket_onl.BookingStateDelete:
		action.actionDelete(ticket)
	case ticket_onl.BookingStateSancelled:
		fallthrough
	case ticket_onl.BookingStateFinished:
		action.actionFinish(ticket)
	case ticket_onl.BookingStateNotArrived:
	case ticket_onl.BookingStateCheckCode:
		action.actionCheckCode()
	case ticket_onl.BookingStateCreateCetm:
		action.actionCreateCetm()
	case ticket_onl.BookingCustomerUpdate:
		action.cusUpdate(ticket)
	default:
		err := errors.New("No Action")
		action.SetError(err)
	}
}
