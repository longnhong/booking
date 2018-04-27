package ticket

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/system"
	"cetm_booking/x/mrw/encode"
	"ehelp/x/rest"
)

func ActionChange(tkID string, cusId string, actionStatus ticket_onl.BookingState, extra encode.RawMessage) *ticket_onl.TicketBooking {
	var action = system.NewTicketAction()
	action.CusID = cusId
	action.TicketID = tkID
	action.Action = actionStatus
	action.Extra = extra
	system.TicketWorkerDay.TriggerTicketAction(action)
	tk, err := action.Wait()
	rest.AssertNil(rest.BadRequestValid(err))
	return tk
}
