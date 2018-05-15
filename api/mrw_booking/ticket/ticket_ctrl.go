package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/system"
	"cetm_booking/x/mrw/encode"
	"ehelp/x/rest"
)

func actionChange(tkID string, cusID string, actionStatus ticket_onl.BookingState, extra encode.RawMessage) *ticket_onl.TicketBooking {
	var action = system.NewTicketAction()
	action.CusID = cusID
	action.TicketID = tkID
	action.Action = actionStatus
	action.Extra = extra
	system.TicketWorkerDay.TriggerTicketAction(action)
	tk, err := action.Wait()
	rest.AssertNil(rest.BadRequestValid(err))
	return tk
}

func setBankTickets(branchID string, serviceID string, timeStart int64, timeEnd int64) (*bankTickets, error) {
	var reslt, err = ticket_onl.GetTicketTimeInBranch(branchID, timeStart, timeEnd)
	if err != nil {
		return nil, err
	}
	var result = make([]resTime, len(reslt))
	for i, item := range reslt {
		var res = resTime{
			ID:         item.ID,
			TimeGoBank: item.TimeGoBank,
			TypeTicket: item.TypeTicket,
			ServiceID:  item.ServiceID,
		}
		result[i] = res
	}
	data, err := ctrl.SearchBank(branchID, serviceID)
	if err != nil {
		return nil, err
	}
	return &bankTickets{Bank: data, Tickets: result}, nil
}
