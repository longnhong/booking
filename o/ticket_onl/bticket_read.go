package ticket_onl

import (
	"cetm_booking/x/math"
	"cetm_booking/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func CheckCustomerCode(customerCode string, branchID string) (tk *TicketBooking, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_code": customerCode,
		"branch_id":     branchID,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		//"status": BookingStateCreated,
	}
	err = TicketBookingTable.FindOne(queryMatch, &tk)
	if err != nil {
		if err.Error() == "not found" {
			err = rest.BadRequestValid(errors.New("Code không tồn tại!"))
		}
		return
	}
	if tk.Status == BookingStateSancelled || tk.Status == BookingStateFinished {
		err = errors.New("Vé đã kết thúc!")
	} else if tk.CheckInAt != 0 {
		err = rest.BadRequestValid(errors.New("Code đã sử dụng!"))
	}
	return
}

func GetCustomerMySchedule(customerId string) (btks []*RateTicket, err error) {
	var status = []string{string(BookingStateSancelled), (string(BookingStateFinished))}
	var queryMatch = bson.M{
		"customer_id": customerId,
		"updated_at":  bson.M{"$ne": 0},
		"status":      bson.M{"$in": status},
	}
	var query = []bson.M{}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "_id",
		"foreignField": "bticket_id",
		"as":           "rate",
	}
	var unWindRate = bson.M{"path": "$rate", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinRate},
		{"$unwind": unWindRate},
	}
	err = TicketBookingTable.Pipe(query).All(&btks)
	rest.IsErrorRecord(err)
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
		"localField":   "_id",
		"foreignField": "bticket_id",
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

func CheckTicketByDay(customerId string) (btks []*TicketBooking, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_id": customerId,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BookingStateCreated,
	}
	err = TicketBookingTable.FindWhere(queryMatch, &btks)
	return btks, rest.IsErrorRecord(err)
}

func GetTicketNear(customerId string) (btk *RateTicket, err error) {
	var queryMatch = bson.M{
		"customer_id": customerId,
		"status":      BookingStateFinished,
	}
	var query = []bson.M{}
	var sort = bson.M{
		"created_at": -1,
	}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "_id",
		"foreignField": "bticket_id",
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
	var timeNow = math.GetTimeNowVietNam().Unix()
	var tracks = tk.updateTrack(tk.ServiceID, tk.BranchID, BookingStateConfirmed, timeNow)
	var up = bson.M{
		"updated_at":  time.Now().Unix(),
		"check_in_at": timeNow,
		"status":      BookingStateConfirmed,
		"tracks":      tracks,
	}
	var err = TicketBookingTable.UnsafeUpdateByID(tk.ID, up)
	if err == nil {
		tk.CheckInAt = timeNow
		tk.Status = BookingStateConfirmed
		tk.Tracks = tracks
	}
	return err
}

func (tk *TicketBooking) UpdateByCnumCetm(cnum string, idCetm string) error {
	var up = bson.M{
		"cnum_cetm":      cnum,
		"id_ticket_cetm": idCetm,
		"status":         BookingStateConfirmed,
	}
	var err = TicketBookingTable.UnsafeUpdateByID(tk.ID, up)
	if err == nil {
		tk.CnumCetm = cnum
		tk.IdTicketCetm = idCetm
		tk.Status = BookingStateConfirmed
	}
	return err
}

func GetTicketDayInBranch(branchID string, timeStart int64, timeEnd int64) (btks []*TicketUser, err error) {
	var queryMatch = bson.M{
		"branch_id": branchID,
		"time_go_bank": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
		"status": BookingStateCreated,
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
	rest.IsErrorRecord(err)
	return btks, err
}
func GetTicketTimeInBranch(branchID string, timeStart int64, timeEnd int64) (btks []*TicketBooking, err error) {
	var queryMatch = bson.M{
		"branch_id": branchID,
		"time_go_bank": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
		"status": BookingStateCreated,
	}

	err = TicketBookingTable.FindWhere(queryMatch, &btks)
	return btks, err
}

func GetAllTicketDay() (btks []*TicketBooking, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BookingStateCreated,
	}
	return btks, TicketBookingTable.FindWhere(queryMatch, &btks)
}

func GetTicketDayByUser(cusId string) (btks []*TicketBooking, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_id": cusId,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BookingStateCreated,
	}
	return btks, TicketBookingTable.FindWhere(queryMatch, &btks)
}

func GetAllTicketByTimeSearch(timeSearch int64, typeTicket TypeTicket) (btks []*TicketBooking, err error) {
	var start, end = math.BeginAndEndDay(timeSearch)
	var queryMatch = bson.M{
		"time_go_bank": bson.M{
			"$gte": start,
			"$lte": end,
		},
		"status": BookingStateCreated,
	}
	if typeTicket == TYPE_SCHEDULE {
		queryMatch["type_ticket"] = TYPE_SCHEDULE
	}
	return btks, TicketBookingTable.FindWhere(queryMatch, &btks)
}

func SearchTicket(idBranchs []string, timeStart int64, timeEnd int64) (btks []*TicketSchedule, err error) {
	//var start, end = math.BeginAndEndDay(timeSearch)
	var queryMatch = bson.M{
		"branch_id": bson.M{"$in": idBranchs},
		"time_go_bank": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
		"status": BookingStateCreated,
	}
	var group = bson.M{"_id": "$branch_id", "count": bson.M{"$sum": 1}}
	var query = []bson.M{
		{"$match": queryMatch},
		{"$group": group},
	}
	return btks, TicketBookingTable.Pipe(query).All(&btks)
}

func GetByID(id string) (tk *TicketBooking, err error) {
	err = TicketBookingTable.FindByID(id, &tk)
	rest.IsErrorRecord(err)
	return
}

func GetTicketByUserNeedFeedback(userId string) (tks []*TicketBooking, err error) {
	var queryMatch = bson.M{
		"customer_id": userId,
		"status":      BookingStateFinished,
	}
	err = TicketBookingTable.FindWhere(queryMatch, &tks)
	rest.IsErrorRecord(err)
	return
}

func UpdateRate(id string, numRate TypeRate) (err error) {
	var up = bson.M{
		"is_rate": numRate,
	}
	err = TicketBookingTable.UnsafeUpdateByID(id, up)
	rest.IsErrorRecord(err)
	return
}

func GetTicketReport(branchIds []string, timeStart int64, timeEnd int64) (btks []*TicketUser, err error) {
	var queryMatch = bson.M{
		"branch_id": bson.M{"$in": branchIds},
		"created_at": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
	}
	var query = []bson.M{}
	var joinUser = bson.M{
		"from":         "user",
		"localField":   "customer_id",
		"foreignField": "_id",
		"as":           "customer",
	}
	var unWindCus = bson.M{"path": "$customer", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinUser},
		{"$unwind": unWindCus},
	}
	err = TicketBookingTable.Pipe(query).All(&btks)
	return btks, err
}
