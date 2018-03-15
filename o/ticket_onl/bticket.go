package ticket_onl

import (
	"cetm_booking/o/auth/user"
	"cetm_booking/x/db/mongodb"
)

var TicketBookingTable = mongodb.NewTable("tk_booking", "tkbk", 18)

type TicketBooking struct {
	mongodb.BaseModel `bson:",inline"`
	TimeGoBank        int64        `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID         string       `json:"service_id" bson:"service_id"`
	BranchID          string       `json:"branch_id"  bson:"branch_id"`
	TypeTicket        TypeTicket   `json:"type_ticket" bson:"type_ticket"`
	Lang              string       `json:"lang" bson:"lang"`
	CustomerCode      string       `bson:"customer_code" json:"customer_code"`
	CustomerID        string       `json:"customer_id"  bson:"customer_id"`
	CheckInAt         int64        `json:"check_in_at"  bson:"check_in_at"`
	IdTicketCetm      string       `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm          string       `json:"cnum_cetm"  bson:"cnum_cetm"`
	TellerID          string       `json:"teller_id"  bson:"teller_id"`
	Teller            string       `json:"teller"  bson:"teller"`
	ServingTime       string       `json:"serving_time"  bson:"serving_time"`
	WatingTime        string       `json:"wating_time"  bson:"wating_time"`
	Status            BookingState `json:"status"  bson:"status"`
}

type UpdateCetm struct {
	BTicketID    string       `bson:"bticket_id" json:"bticket_id"`
	Teller       string       `json:"teller"  bson:"teller"`
	TellerID     string       `json:"teller_id"  bson:"teller_id"`
	IdTicketCetm string       `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm     string       `json:"cnum_cetm"  bson:"cnum_cetm"`
	Status       BookingState `json:"status"  bson:"status"`
	ServingTime  string       `json:"serving_time"  bson:"serving_time"`
	WatingTime   string       `json:"wating_time"  bson:"wating_time"`
}

type WhenCreateTicketInCetm struct {
	CheckInAt    int64  `json:"check_in_at"  bson:"check_in_at"`
	IdTicketCetm string `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm     string `json:"cnum_cetm"  bson:"cnum_cetm"`
}

type TicketUpdate struct {
	UpdatedAt  int64      `json:"updated_at" bson:"updated_at"`
	BTicketID  string     `bson:"bticket_id" json:"bticket_id"`
	TimeGoBank int64      `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID  string     `json:"service_id" bson:"service_id"`
	BranchID   string     `json:"branch_id"  bson:"branch_id"`
	TypeTicket TypeTicket `json:"type_ticket" bson:"type_ticket"`
}

type TicketBookingCreate struct {
	Customer   string     `bson:"customer" json:"customer"`
	TimeGoBank int64      `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID  string     `json:"service_id" bson:"service_id"`
	BranchID   string     `json:"branch_id"  bson:"branch_id"`
	Lang       string     `json:"lang" bson:"lang"`
	CustomerID string     `json:"customer_id" bson:"customer_id"`
	TypeTicket TypeTicket `json:"type_ticket" bson:"type_ticket"`
}

type BookingState string

type TypeTicket string

const (
	TYPE_NOW      = TypeTicket("book_now")
	TYPE_SCHEDUCE = TypeTicket("book_schedule")
)

const (
	BOOKING_STATE_CREATED   = BookingState("created")
	BOOKING_STATE_CONFIRMED = BookingState("confirmed")
	BOOKING_STATE_CANCELLED = BookingState("cancelled")
	BOOKING_STATE_FINISHED  = BookingState("finished")
)

type TicketBranches struct {
	BranchID       string
	TicketBookings []*TicketBooking
}

type TicketUser struct {
	TicketBooking `bson:",inline"`
	Customer      *user.User `json:"customer" bson:"customer"`
}
