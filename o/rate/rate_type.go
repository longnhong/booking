package rate

import (
	"cetm_booking/x/db/mongodb"
)

type Rate struct {
	mongodb.BaseModel `bson:",inline"`
	TicketIdBk        string `bson:"ticket_id" json:"ticket_id"`
	RatePoint         int    `bson:"rate_point" json:"rate_point"`
	Comment           string `bson:"comment" json:"comment"`
	CustomerId        string `bson:"customer_id" json:"customer_id"`
}

var RateTable = mongodb.NewTable("rate", "rt", 20)
