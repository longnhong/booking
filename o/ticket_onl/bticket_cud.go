package ticket_onl

import (
	"cetm_booking/x/math"
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
	err := tkit.updateBf(tkbk)
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
	tkit.BranchName = tkbk.BranchName
	tkit.Tracks = tkbk.Tracks
	return tkit, nil
}

func (tk *TicketBooking) UpdateTicketBookingByCetm(upC *UpdateCetm) error {
	var timeNow = math.GetTimeNowVietNam().Unix()
	upC.Tracks = tk.updateTrack(tk.ServiceID, tk.BranchID, upC.Status, timeNow)
	err := TicketBookingTable.UnsafeUpdateByID(upC.BTicketID, upC)
	if err == nil {
		tk.Tracks = upC.Tracks
		tk.Status = upC.Status
		tk.AvatarTeller = upC.AvatarTeller
		tk.TellerID = upC.TellerID
		tk.IdTicketCetm = upC.IdTicketCetm
		tk.CnumCetm = upC.CnumCetm
		tk.ServingTime = upC.ServingTime
		tk.WaitingTime = upC.WaitingTime
		tk.Teller = upC.Teller
	}
	return err
}

func (tk *TicketBooking) MarkDeleteTicket() error {
	var timeNow = math.GetTimeNowVietNam().Unix()
	var tracks = tk.updateTrack(tk.ServiceID, tk.BranchID, BOOKING_STATE_DELETE, timeNow)
	var updateCancel = bson.M{
		"status":     BOOKING_STATE_DELETE,
		"updated_at": 0,
		"tracks":     tracks,
	}
	err := TicketBookingTable.UnsafeUpdateByID(tk.ID, updateCancel)
	if err == nil {
		tk.Tracks = tracks
	}
	return err
}

func UpdateStatusTickets(ids []string, status BookingState) (error, int) {
	var newUp = map[string]interface{}{
		"status": status,
	}
	var rest, err = TicketBookingTable.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": newUp})
	return err, rest.Updated
}

func UpdateMissedTickets() (error, int) {
	var timeNow = math.GetTimeNowVietNam().Unix()
	var queryMatch = bson.M{
		"time_go_bank": bson.M{
			"$lte": timeNow,
		},
		"status": BOOKING_STATE_CREATED,
	}
	var newUp = map[string]interface{}{
		"status":     BOOKING_STATE_NOT_ARRIVED,
		"updated_at": timeNow,
	}
	var rest, err = TicketBookingTable.UpdateAll(queryMatch, bson.M{"$set": newUp})
	return err, rest.Updated
}
