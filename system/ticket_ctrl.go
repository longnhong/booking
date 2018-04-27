package system

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/utils"
	"fmt"
)

func (c *ticketWorker) TicketWorking(action *TicketAction) error {
	if action == nil {
		return nil
	}
	defer utils.Recover()
	defer action.Done()
	action.handlerAction()
	var err = action.GetError()
	if err == nil {
		var timegoBank = action.Ticket.TimeGoBank
		if action.Action == ticket_onl.BOOKING_STATE_CREATED && math.CompareDayTime(math.GetTimeNowVietNam(), timegoBank) == 0 {
			var hourDay = math.HourMinuteEpoch(timegoBank)
			fmt.Printf("TẠO VÉ TRONG NGÀY", hourDay)
			var tkDay = ticket_onl.TicketDay{
				TicketBooking: action.Ticket,
				HourTimeGo:    hourDay,
			}
			c.TicketCaches[action.Ticket.ID] = &tkDay
		}
	}
	return err
}
