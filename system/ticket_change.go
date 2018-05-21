package system

import (
	"cetm_booking/o/ticket_onl"
	"errors"
)

func (action *TicketAction) handlerAction(ticket *ticket_onl.TicketBooking) {
	switch action.Action {
	case ticket_onl.BOOKING_STATE_CREATED:
		action.actionCreate()
	//case ticket_onl.BOOKING_STATE_CONFIRMED:
	case ticket_onl.BOOKING_STATE_DELETE:
		action.actionDelete(ticket)
	case ticket_onl.BOOKING_STATE_CANCELLED:
		//fallthrough
	case ticket_onl.BOOKING_STATE_FINISHED:
		action.actionFinish(ticket)
	case ticket_onl.BOOKING_STATE_NOT_ARRIVED:
	case ticket_onl.BOOKING_STATE_CHECK_CODE:
		action.actionCheckCode()
	case ticket_onl.BOOKING_STATE_CREATE_CETM:
		action.actionCreateCetm()
	case ticket_onl.BOOKING_CUSTOMER_UPDATE:
		action.cusUpdate(ticket)
	default:
		err := errors.New("No Action")
		action.SetError(err)
	}
}
