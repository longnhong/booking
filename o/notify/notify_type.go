package notify

import (
	"cetm_booking/x/db/mongodb"
)

var NotifyTable = mongodb.NewTable("notify", "k", 20)

type Notify struct {
	mongodb.BaseModel `bson:",inline"`
	Title             string `bson:"title" json:"title" validate:"required"`
	Description       string `bson:"description" json:"description" validate:"required"`
	CustomerId        string `bson:"customer_id" json:"customer_id" validate:"required"`
	URL               string `bson:"url" json:"url"`
	IsReaded          bool   `bson:"is_readed" json:"is_readed"`
	BticketID         string `bson:"bticket_id" json:"bticket_id"`
}
