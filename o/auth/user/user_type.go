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
}

type Role int

var UserTable = mongodb.NewTable("user", "usr", 20)

var ROLE_CETM = Role(1)
var ROLE_USER = Role(2)
