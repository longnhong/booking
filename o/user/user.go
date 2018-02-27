package user

import (
	"cetm_booking/x/db/mongodb"
)

type User struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string   `bson:"name" json:"name" validate:"required"`
	UserName          string   `bson:"uname" json:"uname" validate:"required"`
	Password          Password `bson:"password" json:"password" validate:"required"`
	Role              Role     `bson:"role" json:"role"`
}

type Owner struct {
	User    `bson:",inline"`
	Address string `bson:"address" json:"address" validate:"required"`
	Area    string `bson:"area" json:"area" validate:"required"`
	Phone   string `bson:"phone" json:"phone" validate:"required"`
}

type Staff struct {
	Owner          `bson:",inline"`
	BirthDate      string   `bson:"birth_date" json:"birth_date" `
	IdentityNumber string   `bson:"identity_number" json:"identity_number" validate:"required"`
	Service        []string `bson:"services" json:"services"`
	Certificate    string   `bson:"certificate" json:"certificate" `
	Status         Status   `bson:"status" json:"status"`
}
type Status string

const (
	APPROVE = Status("approve")
	CANCEL  = Status("cancel")
)

type Role string

const (
	SUPER_ADMIN = Role("super_admin")
	ADMIN       = Role("admin")
	STAFF       = Role("staff")
	OWNER       = Role("owner")
)
