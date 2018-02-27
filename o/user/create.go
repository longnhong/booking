package user

import (
	"cetm_booking/x/db/mongodb"
	"cetm_booking/x/rest"
	"cetm_booking/x/rest/validator"
	"gopkg.in/mgo.v2/bson"
)

var UserTable = mongodb.NewTable("user", "usr", 12)

func (u *Staff) Create() error {
	var queryUnique = bson.M{"uname": u.UserName, "role": u.Role}
	hashed, _ := u.Password.GererateHashedPassword()
	u.Password = hashed
	if u.Role == STAFF {
		rest.AssertNil(rest.WrapBadRequest(validator.Validate(u), ""))
		return UserTable.CreateUnique(queryUnique, u)
	}
	if u.Role == OWNER {
		rest.AssertNil(rest.WrapBadRequest(validator.Validate(u.Owner), ""))
		return UserTable.CreateUnique(queryUnique, &u.Owner)
	}
	{
		rest.AssertNil(rest.WrapBadRequest(validator.Validate(u.User), ""))
		if u.Role == SUPER_ADMIN {
			return UserTable.CreateUnique(bson.M{"role": SUPER_ADMIN}, &u.User)
		}
		return UserTable.CreateUnique(queryUnique, &u.User)
	}
}
