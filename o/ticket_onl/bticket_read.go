package ticket_onl

import (
	"cetm_booking/x/utility"
	"errors"
	"gopkg.in/mgo.v2/bson"
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
		//"status":        BOOKING_STATE_CREATED,
	}
	err = TicketBookingTable.FindOne(queryMatch, &tk)
	if err != nil {
		return
	}
	if tk.Status == BOOKING_STATE_CANCELLED || tk.Status == BOOKING_STATE_FINISHED {
		err = errors.New("Vé đã kết thúc!")
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
	return btks, TicketBookingTable.FindWhere(queryMatch, &btks)
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
