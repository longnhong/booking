package ticket_onl

import (
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

func (tkit *TicketBooking) UpdateTicketBookingByCustomer(tkbk *TicketUpdate) (*TicketBooking, error) {
	err := tkbk.updateBf()
	if err != nil {
		return nil, err
	}
	err = TicketBookingTable.UnsafeUpdateByID(tkbk.BTicketID, tkbk)
	if err != nil {
		return nil, err
	}
	tkit.BranchAddress = tkbk.BranchAddress
	tkit.BranchID = tkbk.BranchID
	tkit.TimeGoBank = tkbk.TimeGoBank
	tkit.ServiceID = tkbk.ServiceID
	tkit.ServiceName = tkbk.ServiceName
	tkit.TypeTicket = tkbk.TypeTicket
	tkit.UpdatedAt = tkbk.UpdatedAt
	return tkit, nil
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

func UpdateStatusTickets(ids []string, status BookingState) (error, int) {
	var newUp = map[string]interface{}{
		"status": status,
	}
	var rest, err = TicketBookingTable.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": newUp})
	return err, rest.Updated
}
