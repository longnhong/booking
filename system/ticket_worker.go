package system

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/mrw/event"
)

type TicketWorker struct {
	TicketCaches map[string]*ticket_onl.TicketDay
	TicketUpdate chan *TicketAction
	doneAction   *event.Hub
}

func (tw *TicketWorker) TriggerTicketAction(action *TicketAction) {
	tw.TicketUpdate <- action
}

func newCacheTicketWorker() *TicketWorker {
	return &TicketWorker{
		TicketCaches: make(map[string]*ticket_onl.TicketDay, 0),
		TicketUpdate: make(chan *TicketAction, event.MediumHub),
	}
}

func (tkDay *TicketWorker) GetTicketByID(idTicket string) (*ticket_onl.TicketBooking, error) {
	if tk, ok := tkDay.TicketCaches[idTicket]; ok {
		return tk.TicketBooking, nil
	}
	ticket, err := ticket_onl.GetByID(idTicket)
	return ticket, err
}

func (tw *TicketWorker) OnActionDone() (event.Line, event.Cancel) {
	return tw.doneAction.NewLine()
}

func (tkDay *TicketWorker) removeTksTicketWorkerDay() {
	if tkDay != nil {
		if len(tkDay.TicketCaches) > 0 {
			for k := range tkDay.TicketCaches {
				delete(tkDay.TicketCaches, k)
			}
		}
	}
}
