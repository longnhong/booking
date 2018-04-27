package ticket_onl

import (
	"cetm_booking/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func (tkbkCreate *TicketBookingCreate) CrateTicketBooking() (*TicketBooking, error) {
	var err, tkbk = tkbkCreate.createBf()
	if err != nil {
		return nil, err
	}
	err = TicketBookingTable.Create(tkbk)
	return tkbk, err
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

func (upC *UpdateCetm) UpdateTicketBookingByCetm() error {
	return TicketBookingTable.UnsafeUpdateByID(upC.BTicketID, upC)
}

func (tk *TicketBooking) MarkDeleteTicket() error {
	var updateCancel = bson.M{
		"status":     BOOKING_STATE_DELETE,
		"updated_at": 0,
	}
	err := TicketBookingTable.UnsafeUpdateByID(tk.ID, updateCancel)
	return err
}

func UpdateStatusTickets(tks []*TicketDay, status BookingState) (error, int) {

	var newUp = map[string]interface{}{
		"status": status,
	}
	var rest, err = TicketBookingTable.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": newUp})
	return err, rest.Updated
}
