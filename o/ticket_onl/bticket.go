package ticket_onl

import (
	"cetm_booking/o/auth/user"
	"cetm_booking/o/notify"
	"cetm_booking/o/rate"
	"cetm_booking/x/db/mongodb"
)

var TicketBookingTable = mongodb.NewTable("tk_booking", "tkbk", 18)

type TicketBooking struct {
	mongodb.BaseModel `bson:",inline"`
	TimeGoBank        int64        `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID         string       `json:"service_id" bson:"service_id"`
	ServiceName       string       `json:"service_name" bson:"service_name"`
	BranchID          string       `json:"branch_id"  bson:"branch_id"`
	BranchAddress     string       `json:"branch_address"  bson:"branch_address"`
	TypeTicket        TypeTicket   `json:"type_ticket" bson:"type_ticket"`
	Lang              string       `json:"lang" bson:"lang"`
	CustomerCode      string       `bson:"customer_code" json:"customer_code"`
	CustomerID        string       `json:"customer_id"  bson:"customer_id"`
	CheckInAt         int64        `json:"check_in_at"  bson:"check_in_at"`
	AvatarTeller      string       `json:"avatar_teller"  bson:"avatar_teller"`
	IdTicketCetm      string       `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	BranchName        string       `json:"branch_name"  bson:"branch_name"`
	Tracks            []TicketHst  `json:"tracks"  bson:"tracks"`
	CnumCetm          string       `json:"cnum_cetm"  bson:"cnum_cetm"`
	TellerID          string       `json:"teller_id"  bson:"teller_id"`
	ServingAt         int64        `json:"serving_at"  bson:"serving_at"`
	Teller            string       `json:"teller"  bson:"teller"`
	ServingTime       int64        `json:"serving_time"  bson:"serving_time"`
	WaitingTime       int64        `json:"waiting_time"  bson:"waiting_time"`
	IsRate            TypeRate     `json:"is_rate"  bson:"is_rate"` //0: chưa rate, 1:rate, 2: khong rate
	Status            BookingState `json:"status"  bson:"status"`
}

type TicketHst struct {
	ServiceID string       `json:"service_id" bson:"service_id"`
	BranchID  string       `json:"branch_id"  bson:"branch_id"`
	MTime     int64        `json:"ctime"  bson:"ctime"`
	Status    BookingState `json:"status"  bson:"status"`
}

type RateTicket struct {
	TicketBooking `bson:",inline"`
	Rate          *rate.Rate `json:"rate"  bson:"rate"`
}

type UpdateCetm struct {
	BTicketID    string       `bson:"bticket_id" json:"bticket_id"`
	Teller       string       `json:"teller"  bson:"teller"`
	AvatarTeller string       `json:"avatar_teller"  bson:"avatar_teller"`
	TellerID     string       `json:"teller_id"  bson:"teller_id"`
	IdTicketCetm string       `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	Tracks       []TicketHst  `json:"tracks"  bson:"tracks"`
	CnumCetm     string       `json:"cnum_cetm"  bson:"cnum_cetm"`
	Status       BookingState `json:"status"  bson:"status"`
	ServingTime  int64        `json:"serving_time"  bson:"serving_time"`
	WaitingTime  int64        `json:"waiting_time"  bson:"waiting_time"`
	ServingAt    int64        `json:"serving_at"  bson:"serving_at"`
}

type WhenCreateTicketInCetm struct {
	CheckInAt    int64  `json:"check_in_at"  bson:"check_in_at"`
	IdTicketCetm string `json:"id_ticket_cetm"  bson:"id_ticket_cetm"`
	CnumCetm     string `json:"cnum_cetm"  bson:"cnum_cetm"`
}

type TicketUpdate struct {
	UpdatedAt     int64        `json:"updated_at" bson:"updated_at"`
	BTicketID     string       `bson:"bticket_id" json:"bticket_id"`
	TimeGoBank    int64        `bson:"time_go_bank" json:"time_go_bank"`
	ServiceID     string       `json:"service_id" bson:"service_id"`
	ServiceName   string       `json:"service_name" bson:"service_name"`
	BranchName    string       `json:"branch_name"  bson:"branch_name"`
	BranchID      string       `json:"branch_id"  bson:"branch_id"`
	BranchAddress string       `json:"branch_address"  bson:"branch_address"`
	TypeTicket    TypeTicket   `json:"type_ticket" bson:"type_ticket"`
	Tracks        []TicketHst  `json:"tracks"  bson:"tracks"`
	Status        BookingState `json:"status"  bson:"status"`
}

type TicketDay struct {
	*TicketBooking `bson:",inline"`
	HourTimeGo     float32 `json:"hour_time_go" bson:"hour_time_go"`
	IsUsedPush     bool
	IsUsedNear     bool
	IsUsedOut      bool
	Customer       *user.User
}

type TicketSchedule struct {
	IdBranch    string `json:"id" bson:"_id"`
	CountPeople int    `json:"count" bson:"count"`
}

type TicketBookingCreate struct {
	Customer      string     `bson:"customer" json:"customer"`
	TimeGoBank    int64      `bson:"time_go_bank" json:"time_go_bank"`
	BranchAddress string     `json:"branch_address"  bson:"branch_address"`
	BranchName    string     `json:"branch_name"  bson:"branch_name"`
	ServiceID     string     `json:"service_id" bson:"service_id"`
	ServiceName   string     `json:"service_name" bson:"service_name"`
	BranchID      string     `json:"branch_id"  bson:"branch_id"`
	Lang          string     `json:"lang" bson:"lang"`
	CustomerID    string     `json:"customer_id" bson:"customer_id"`
	TypeTicket    TypeTicket `json:"type_ticket" bson:"type_ticket"`
}

type BookingState string

type TypeTicket string

type TypeRate int

const (
	TypeDefaultRate = TypeRate(0)
	TypeRated       = TypeRate(1)
	TypeNoRate      = TypeRate(2)
)

const (
	TypeNow      = TypeTicket("book_now")
	TypeSchedule = TypeTicket("book_schedule")
)

const (
	BookingStateCreated    = BookingState("created")     // vừa tạo
	BookingStateConfirmed  = BookingState("confirmed")   // đã xác nhận đến ngân hàng
	BookingStateSancelled  = BookingState("cancelled")   // cetm gọi ko có mặt
	BookingStateDelete     = BookingState("deleted")     // kh xóa vé
	BookingStateFinished   = BookingState("finished")    // đã giao dịch xong
	BookingStateNotArrived = BookingState("not_arrived") // không đến
	BookingStateCheckCode  = BookingState("check_code")  // CheckCode
	BookingStateCreateCetm = BookingState("create_cetm") // CheckCode
	BookingCustomerUpdate  = BookingState("cus_update")  //kh update
)

type TicketBranches struct {
	BranchID       string
	TicketBookings []*TicketBooking
}

type TicketUser struct {
	TicketBooking `bson:",inline"`
	Customer      *user.User `json:"customer" bson:"customer"`
}

type NotifyTicket struct {
	Notify *notify.Notify `json:"notify"`
	Ticket *TicketBooking `json:"ticket_booking"`
}
