package ticket_onl

import (
	"cetm_booking/x/utility"
	"gopkg.in/mgo.v2/bson"
)

func CheckCustomerCode(customerCode string, branchID string) (int, error) {
	var queryMatch = bson.M{
		"customer_Code": customerCode,
		"branch_id":     branchID,
		"status":        BOOKING_STATE_CREATED,
	}
	return TicketBookingTable.CountWhere(queryMatch)
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
