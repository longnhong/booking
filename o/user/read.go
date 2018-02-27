package user

import (
	"gopkg.in/mgo.v2/bson"
)

func GetByUNamePwd(uname, pwd, role string) (*User, error) {
	var user *User
	var query = bson.M{"uname": uname, "role": role}
	if role == string(ADMIN) {
		query = bson.M{
			"uname": uname,
			"$or": []bson.M{
				bson.M{"role": ADMIN},
				bson.M{"role": SUPER_ADMIN},
			},
		}
	}
	err := UserTable.FindOne(query, &user)
	if err != nil {
		return nil, err
	}
	if err := user.Password.ComparePassword(pwd); err == nil {
		return user, nil
	}
	return user, nil
}

func GetSuperUser() *User {
	var user *User
	err := UserTable.FindOne(bson.M{"role": SUPER_ADMIN}, &user)
	if err != nil {
		return nil
	}
	return user
}
