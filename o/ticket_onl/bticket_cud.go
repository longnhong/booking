package ticket_onl

import (
	"cetm_booking/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func (tkbkCreate *TicketBookingCreate) CrateTicketBooking() *TicketBooking {
	var err, tkbk = tkbkCreate.createBf()
	rest.AssertNil(err)
	rest.AssertNil(TicketBookingTable.Create(tkbk))
	return tkbk
}

func (tkbk *TicketUpdate) UpdateTicketBookingByCustomer() (tk *TicketBooking) {
	var err = TicketBookingTable.FindByID(tkbk.BTicketID, &tk)
	if err != nil || tk == nil {
		rest.AssertNil(errors.New("Không tồn tại vé này!"))
	}
	rest.AssertNil(tkbk.updateBf())
	rest.AssertNil(TicketBookingTable.UnsafeUpdateByID(tkbk.BTicketID, tkbk))
	tk.TimeGoBank = tkbk.TimeGoBank
	tk.ServiceID = tkbk.ServiceID
	tk.BranchID = tkbk.BranchID
	tk.TypeTicket = tkbk.TypeTicket
	tk.UpdatedAt = tkbk.UpdatedAt
	return
}

func (upC *UpdateCetm) UpdateTicketBookingByCetm() {
	rest.AssertNil(TicketBookingTable.UnsafeUpdateByID(upC.BTicketID, upC))
}

func MarkDeleteTicket(id string) (tk *TicketBooking) {
	var updateCancel = bson.M{
		"status":     BOOKING_STATE_DELETE,
		"updated_at": 0,
	}
	rest.AssertNil(TicketBookingTable.UnsafeUpdateByID(id, updateCancel))
	var err error
	tk, err = GetByID(id)
	rest.AssertNil(err)
	return tk
}
