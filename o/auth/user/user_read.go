package user

import (
	"cetm_booking/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func GetUserByLogin(username string, password string) (usr *User, err error) {
	err = UserTable.FindOne(bson.M{"username": username}, &usr)
	if err != nil || usr == nil {
		return nil, rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!"))
	}
	if err := usr.Password.ComparePassword(password); err != nil {
		return nil, rest.BadRequestValid(errors.New("Password sai!"))
	}
	return usr, nil
}

func GetByID(idUser string) (usr *User, err error) {
	return usr, UserTable.FindByID(idUser, &usr)
}
