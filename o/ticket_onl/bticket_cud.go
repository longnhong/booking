package ticket_onl

import (
	"cetm_booking/x/rest"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (tkbkCreate *TicketBookingCreate) CrateTicketBooking() *TicketBooking {
	var err, tkbk = tkbkCreate.createBf()
	rest.AssertNil(err)
	rest.AssertNil(TicketBookingTable.Create(tkbk))
	return tkbk
}

func (tkbk *TicketUpdate) UpdateTicketBookingByCustomer() {
	rest.AssertNil(tkbk.updateBf())
	rest.AssertNil(TicketBookingTable.UnsafeUpdateByID(tkbk.BTicketID, tkbk))
}

func (upC *UpdateCetm) UpdateTicketBookingByCetm() {
	rest.AssertNil(TicketBookingTable.UnsafeUpdateByID(upC.BTicketID, upC))
}

func CancleTicket(id string) {
	var updateCancel = bson.M{
		"status": BookingStateCancelled,
		"dtime":  time.Now().Unix(),
	}
	rest.AssertNil(TicketBookingTable.UnsafeUpdateByID(id, updateCancel))
}
