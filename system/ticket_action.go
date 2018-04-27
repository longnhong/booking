package system

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/mlog"
	"cetm_booking/x/mrw/encode"
)

var logAction = mlog.NewTagLog("system")

type TicketAction struct {
	Action    ticket_onl.BookingState
	TicketID  string
	CusID     string
	PushToken string
	Extra     encode.RawMessage         `json:"extra"`
	Ticket    *ticket_onl.TicketBooking `json:"ticket"`
	Error     TicketActionError         `json:"error"`
	doneC     chan struct{}
	used      bool // must be trigger at most once
}
type TicketActionError struct {
	s string
}

func (a *TicketAction) Done() bool {
	a.doneC <- struct{}{}
	return a.GetError() == nil
}

func (e *TicketActionError) Error() string {
	return e.s
}

func (e *TicketActionError) GetError() error {
	if len(e.s) > 0 {
		return e
	}
	return nil
}

func (a *TicketAction) SetError(err error) {
	if err == nil {
		return
	}
	logAction.Errorf("SetError", err)
	a.Error = TicketActionError{s: err.Error()}
}

func (a *TicketAction) GetError() error {
	return a.Error.GetError()
}

func (a *TicketAction) Wait() (*ticket_onl.TicketBooking, error) {
	if a.doneC == nil {
		logAction.Errorf("Wait()", "no done channel")
		panic("no done channel")
	}
	<-a.doneC
	return a.Ticket, a.GetError()
}

func NewTicketAction() *TicketAction {
	var a = &TicketAction{
		doneC: make(chan struct{}, 1),
	}
	return a
}
