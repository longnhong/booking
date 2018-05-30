package system

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/utils"
)

func (c *TicketWorker) TicketWorking(action *TicketAction) error {
	if action == nil {
		return nil
	}
	defer utils.Recover()
	defer action.Done()
	var ticket *ticket_onl.TicketBooking
	var err error
	var tkAction = action.Action
	if checkStatusGetTicket(tkAction) {
		ticket, err = c.GetTicketByID(action.TicketID)
		if err != nil {
			action.SetError(err)
			return err
		}
	}
	action.handlerAction(ticket)
	err = action.GetError()
	if err == nil {
		var timegoBank = action.Ticket.TimeGoBank
		if action.Action == ticket_onl.BookingStateCreated && math.CompareDayTime(math.GetTimeNowVietNam(), timegoBank) == 0 {
			var hourDay = math.HourMinuteEpoch(timegoBank)
			var tkDay = ticket_onl.TicketDay{
				TicketBooking: action.Ticket,
				HourTimeGo:    hourDay,
			}
			c.TicketCaches[action.Ticket.ID] = &tkDay
		} else if tkAction == ticket_onl.BookingStateCheckCode {
			if val, ok := c.TicketCaches[action.Ticket.ID]; ok {
				val.TicketBooking = action.Ticket
			}
		}
	}
	return err
}

func checkStatusGetTicket(tkAction ticket_onl.BookingState) bool {
	switch tkAction {
	case ticket_onl.BookingStateCreated:
		return false
	case ticket_onl.BookingStateCheckCode:
		return false
	case ticket_onl.BookingStateCreateCetm:
		return false
	default:
		return true
	}
}
