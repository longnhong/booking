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

func SetBankTickets(branchId string, serviceID string, timeStart int64, timeEnd int64) (bak *bankTickets) {
	var reslt, err = ticket_onl.GetTicketTimeInBranch(branchId, timeStart, timeEnd)

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
	rest.AssertNil(err)
	var data = SearchBank(branchId, serviceID)
	bak = &bankTickets{
		Bank:    data,
		Tickets: result,
	}
	return
}
