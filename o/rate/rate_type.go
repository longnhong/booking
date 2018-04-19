package rate

import (
	"cetm_booking/x/db/mongodb"
)

type Rate struct {
	mongodb.BaseModel `bson:",inline"`
	TicketIdBk        string  `bson:"bticket_id" json:"bticket_id"`
	RatePoint         float32 `bson:"rate_point" json:"rate_point"`
	Comment           string  `bson:"comment" json:"comment"`
	CustomerId        string  `bson:"customer_id" json:"customer_id"`
}

var RateTable = mongodb.NewTable("rate", "rt", 20)
