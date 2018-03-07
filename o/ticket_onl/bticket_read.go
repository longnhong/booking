package ticket_onl

import (
	"cetm_booking/x/utility"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func CheckCustomerCode(customerCode string, branchID string) (tk *TicketBooking, err error) {
	var queryMatch = bson.M{
		"customer_code": customerCode,
		"branch_id":     branchID,
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

func CheckCustomerIdByDay(customerId string) (*TicketBooking, error) {
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
	var ticket *TicketBooking
	return ticket, TicketBookingTable.FindOne(queryMatch, &ticket)
}
