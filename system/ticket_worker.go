package system

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/mrw/event"
	"errors"
)

type ticketWorker struct {
	TicketCaches map[string]*ticket_onl.TicketDay
	TicketUpdate chan *TicketAction
	doneAction   *event.Hub
}

func (tw *ticketWorker) TriggerTicketAction(action *TicketAction) {
	tw.TicketUpdate <- action
}

var TicketWorkerDay = newCacheTicketWorker()

func newCacheTicketWorker() *ticketWorker {
	return &ticketWorker{
		TicketCaches: make(map[string]*ticket_onl.TicketDay, 0),
		TicketUpdate: make(chan *TicketAction, event.MediumHub),
	}
}

func GetTicketByID(idTicket string) (*ticket_onl.TicketBooking, error) {
	if tk, ok := TicketWorkerDay.TicketCaches[idTicket]; ok {
		return tk.TicketBooking, nil
	}
	ticket, err := ticket_onl.GetByID(idTicket)
	return ticket, err
}

func CheckCode(branchID string, code string) (tk *ticket_onl.TicketBooking, err error) {
	var cacheTickets = TicketWorkerDay.TicketCaches
	for _, tkDay := range cacheTickets {
		if tkDay.BranchID == branchID && code == tkDay.CustomerCode {
			tk = tkDay.TicketBooking
			var resTime = math.HourMinute() - math.HourMinuteEpoch(tk.TimeGoBank)
			if (resTime < 0.15 && resTime > 0) || (resTime < 0 && resTime > 0.15) {
				return
			} else {
				err = errors.New("Chưa đến giờ!")
			}
		}
	}
	return
}

func (tw *ticketWorker) OnActionDone() (event.Line, event.Cancel) {
	return tw.doneAction.NewLine()
}
