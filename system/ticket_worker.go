package system

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/mrw/event"
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

func (tw *ticketWorker) OnActionDone() (event.Line, event.Cancel) {
	return tw.doneAction.NewLine()
}

func RemoveTksTicketWorkerDay() {
	if TicketWorkerDay != nil && len(TicketWorkerDay.TicketCaches) > 0 {
		for k := range TicketWorkerDay.TicketCaches {
			delete(TicketWorkerDay.TicketCaches, k)
		}
	}
}
