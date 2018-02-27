package user

import (
	"cetm_booking/x/db/mongodb"
)

type UserInterface struct {
	mongodb.BaseModel `bson:",inline"`
	FullName          string  `bson:"full_name" json:"full_name" validate:"required"`
	Password          string  `bson:"password" json:"password,omitempty"`
	COD               float32 `bson:"cod" json:"cod"`
	Phone             string  `bson:"phone" json:"phone" validate:"required"`
	Verified          bool    `bson:"verified" json:"verified"`
	DateOfBirth       string  `bson:"date_of_birth" json:"date_of_birth"`
	Address           string  `bson:"address" json:"address"`
	// VehicleName         string  `bson:"vehicle_name" json:"vehicle_name"`
	// VehicleType         string  `bson:"vehicle_type" json:"vehicle_type"`
	// VehicleRegistration string  `bson:"vehicle_registration" json:"vehicle_registration"`
	// VehicleNumber       string  `bson:"vehicle_number" json:"vehicle_number"`
	Image         string `bson:"image" json:"image"`
	IdentityCard  string `bson:"identity_card" json:"identity_card"`
	CreateDate    int64  `bson:"create_date" json:"create_date"`
	InviteCode    string `bson:"invite_code" json:"invite_code"`
	ReferenceCode string `bson:"reference_code" json:"reference_code"`
	IsActive      bool   `bson:"is_active" json:"is_active"`
	IsVip         bool   `bson:"is_vip" json:"is_vip"`
	Rate5         int    `bson:"rate5" json:"rate5"`
	Rate4         int    `bson:"rate4" json:"rate4"`
	Rate3         int    `bson:"rate3" json:"rate3"`
	Rate2         int    `bson:"rate2" json:"rate2"`
	Rate1         int    `bson:"rate1" json:"rate1"`
	FbID          string `bson:"fb_id" json:"fb_id"`
	FbToken       string `bson:"fb_token" json:"fb_token"`
	GmId          string `bson:"gm_id" json:"gm_id"`
	GmToken       string `bson:"gm_token" json:"gm_token"`
	TypePayment   int    `bson:"type_payment" json:"type_payment"`
}
