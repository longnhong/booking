package user

import (
	"cetm_booking/x/db/mongodb"
)

type User struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string   `bson:"name" json:"name" validate:"required"`
	UserName          string   `bson:"username" json:"username" validate:"required"`
	Password          Password `bson:"password" json:"password" validate:"required"`
	Role              Role     `bson:"role" json:"role"`
	Email             string   `bson:"email" json:"email" `
	Card              string   `bson:"card" json:"card" `
	DateOfBirth       int64    `bson:"date_of_birth" json:"date_of_birth" `
	FullName          string   `bson:"full_name" json:"full_name" `
	PhoneNumber       string   `bson:"phone_number" json:"phone_number" `
	Nationality       string   `bson:"nationality" json:"nationality" `
	VipCode           string   `bson:"vip_code" json:"vip_code" `
	Segment           string   `bson:"segment" json:"segment" `
}

type Role int

var UserTable = mongodb.NewTable("user", "usr", 20)

var ROLE_CETM = Role(1)
var ROLE_USER = Role(2)
