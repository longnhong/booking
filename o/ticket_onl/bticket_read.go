package ticket_onl

import (
	"cetm_booking/x/rest"
	"cetm_booking/x/utility"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func CheckCustomerCode(customerCode string, branchID string) (tk *TicketBooking, err error) {
	var timeBeginDay = utility.BeginningOfDay().Unix()
	var tiemEnOfday = utility.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_code": customerCode,
		"branch_id":     branchID,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		//"status": BOOKING_STATE_CREATED,
	}
	err = TicketBookingTable.FindOne(queryMatch, &tk)
	if err != nil {
		if err.Error() == "not found" {
			err = rest.BadRequestValid(errors.New("Code không tồn tại!"))
		}
		return
	}
	if tk.Status == BOOKING_STATE_CANCELLED || tk.Status == BOOKING_STATE_FINISHED {
		err = errors.New("Vé đã kết thúc!")
	} else if tk.CheckInAt != 0 {
		err = rest.BadRequestValid(errors.New("Code đã sử dụng!"))
	}
	return
}

func GetCustomerIdByDay(customerId string) (btks []*TicketBooking, err error) {
	var timeBeginDay = utility.BeginningOfDay().Unix()
	var tiemEnOfday = utility.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_id": customerId,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BOOKING_STATE_CREATED,
	}
	err = TicketBookingTable.FindWhere(queryMatch, &btks)
	err = rest.IsErrorRecord(err)
	return btks, err
}

func GetAllTicketCus(customerId string) (btks []*RateTicket, err error) {
	var queryMatch = bson.M{
		"customer_id": customerId,
		"updated_at":  bson.M{"$ne": 0},
	}
	var query = []bson.M{}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "ticket_id",
		"foreignField": "_id",
		"as":           "rate",
	}
	var unWindRate = bson.M{"path": "$rate", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinRate},
		{"$unwind": unWindRate},
	}
	return btks, TicketBookingTable.Pipe(query).All(&btks)
}

func CheckTicketByDay(customerId string) (btks *TicketBooking, err error) {
	var timeBeginDay = utility.BeginningOfDay().Unix()
	var tiemEnOfday = utility.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_id": customerId,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BOOKING_STATE_CREATED,
	}
	return btks, TicketBookingTable.FindOne(queryMatch, &btks)
}

func GetTicketNear(customerId string) (btk *RateTicket, err error) {
	var queryMatch = bson.M{
		"customer_id": customerId,
		"status":      BOOKING_STATE_FINISHED,
	}
	var query = []bson.M{}
	var sort = bson.M{
		"created_at": -1,
	}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "ticket_id",
		"foreignField": "_id",
		"as":           "rate",
	}
	var unWindRate = bson.M{"path": "$rate", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinRate},
		{"$unwind": unWindRate},
		{"$sort": sort},
	}

	var btks []*RateTicket
	err = TicketBookingTable.Pipe(query).All(&btks)
	if err == nil && len(btks) > 0 {
		btk = btks[0]
	}
	rest.IsErrorRecord(err)
	return btk, err
}

func (tk *TicketBooking) UpdateTimeCheckIn() error {
	var timeNow = time.Now().Unix()
	var up = bson.M{
		"updated_at":  time.Now().Unix(),
		"check_in_at": timeNow,
		"status":      BOOKING_STATE_CONFIRMED,
	}
	var err = TicketBookingTable.UnsafeUpdateByID(tk.ID, up)
	if err == nil {
		tk.CheckInAt = timeNow
	}
	return err
}

func (tk *TicketBooking) UpdateByCnumCetm(cnum string, idCetm string) error {
	var up = bson.M{
		"cnum_cetm":      cnum,
		"id_ticket_cetm": idCetm,
		"status":         BOOKING_STATE_CONFIRMED,
	}
	var err = TicketBookingTable.UnsafeUpdateByID(tk.ID, up)
	if err == nil {
		tk.CnumCetm = cnum
		tk.IdTicketCetm = idCetm
		tk.Status = BOOKING_STATE_CONFIRMED
	}
	return err
}

func GetTicketDayInBranch(branchID string) (btks []*TicketUser, err error) {
	var timeBeginDay = utility.BeginningOfDay().Unix()
	var tiemEnOfday = utility.EndOfDay().Unix()
	var queryMatch = bson.M{
		"branch_id": branchID,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BOOKING_STATE_CREATED,
	}
	var query = []bson.M{}
	var joinUser = bson.M{
		"from":         "user",
		"localField":   "customer_id",
		"foreignField": "_id",
		"as":           "customer",
	}
	var unWindCus = "$customer"
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinUser},
		{"$unwind": unWindCus},
	}
	err = TicketBookingTable.Pipe(query).All(&btks)
	return btks, err
}

func GetAllTicketDay() (btks []*TicketBooking, err error) {
	var timeBeginDay = utility.BeginningOfDay().Unix()
	var tiemEnOfday = utility.EndOfDay().Unix()
	var queryMatch = bson.M{
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BOOKING_STATE_CREATED,
	}
	return btks, TicketBookingTable.FindWhere(queryMatch, &btks)
}
