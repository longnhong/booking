package user

import (
	//"cetm_booking/o/push_token"
	"cetm_booking/x/rest"
	"errors"
	//	"g/x/web"
	"gopkg.in/mgo.v2/bson"
	//"net/http"
)

func GetUserByLogin(username string, password string) (error, *User) {
	var usr *User
	err := UserTable.FindOne(bson.M{"username": username}, &usr)
	if err != nil || usr == nil {
		return rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!")), nil
	}
	if err := usr.Password.ComparePassword(password); err != nil {
		return rest.BadRequestValid(errors.New("Password sai!")), nil
	}
	return nil, usr
}

type UserToken struct {
	User      `bson:"user" json:"user"`
	RoleToken Role `bson:"role" json:"role"`
}

// func GetUserFromToken(r *http.Request) *UserToken {
// 	var token = web.GetToken(r)
// 	if len(token) < 8 {
// 		panic(web.Unauthorized("missing or invalid access token"))
// 	}
// 	var query []bson.M
// 	var queryMatch = bson.M{}
// 	queryMatch["_id"] = token
// 	var joinCus = bson.M{
// 		"from":         "user",
// 		"localField":   "user_id",
// 		"foreignField": "_id",
// 		"as":           "user",
// 	}
// 	var unWindCus = bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}
// 	query = []bson.M{
// 		{"$match": queryMatch},
// 		{"$lookup": joinCus},
// 		{"$unwind": unWindCus},
// 	}
// 	var usrTk *UserToken
// 	var err = push_token.PushTokenTable.Pipe(query).One(&usrTk)
// 	rest.AssertNil(err)
// 	if usrTk == nil {
// 		rest.AssertNil(web.Unauthorized("missing or invalid access token"))
// 	}
// 	return usrTk
// }
