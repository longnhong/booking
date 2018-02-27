package ticket_onl

import (
	"cetm_booking/x/db/mongodb"
)

var TicketBookingTable = mongodb.NewTable("ticket_booking", "tkbk", 18)

type TicketBooking struct {
	mongodb.BaseModel `bson:",inline"`
	Customer          string       `bson:"customer" json:"customer"`
	TimeGoBank        int64        `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID         string       `json:"service_id" bson:"service_id"`
	BranchID          string       `json:"branch_id"  bson:"branch_id"`
	Lang              string       `json:"lang" bson:"lang"`
	CustomerCode      string       `json:"customer_code"`
	CheckInAt         int64        `json:"check_in_at"  bson:"check_in_at"`
	IdTicketCetm      string       `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm          string       `json:"cnum_cetm"  bson:"cnum_cetm"`
	Status            BookingState `json:"status"  bson:"status"`
}

type UpdateCetm struct {
	BTicketID    string       `bson:"bticket_id" json:"bticket_id"`
	CheckInAt    int64        `json:"check_in_at"  bson:"check_in_at"`
	IdTicketCetm string       `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm     string       `json:"cnum_cetm"  bson:"cnum_cetm"`
	Status       BookingState `json:"status"  bson:"status"`
}

type WhenCreateTicketInCetm struct {
	CheckInAt    int64  `json:"check_in_at"  bson:"check_in_at"`
	IdTicketCetm string `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm     string `json:"cnum_cetm"  bson:"cnum_cetm"`
}

type TicketUpdate struct {
	BTicketID           string `bson:"bticket_id" json:"bticket_id"`
	TicketBookingCreate `bson:",inline"`
}

type TicketBookingCreate struct {
	Customer   string `bson:"customer" json:"customer"`
	TimeGoBank int64  `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID  string `json:"service_id" bson:"service_id"`
	BranchID   string `json:"branch_id"  bson:"branch_id"`
	Lang       string `json:"lang" bson:"lang"`
}

type BookingState string

const (
	BookingStateCreated   = BookingState("created")
	BookingStateCancelled = BookingState("cancelled")
	BookingStateConfirmed = BookingState("confirmed")
	BookingStateUsed      = BookingState("used")
)
