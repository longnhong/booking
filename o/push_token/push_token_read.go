package push_token

import (
	"cetm_booking/x/rest"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*PushToken, error) {
	var auth *PushToken
	return auth, PushTokenTable.FindOne(bson.M{
		"_id":       id,
		"is_revoke": false,
	}, &auth)
}

func GetPushsUserId(userId string) ([]string, error) {
	var pushs []string
	var err = PushTokenTable.Find(bson.M{
		"user_id": userId,
	}).Distinct("push_token", &pushs)
	fmt.Println(len(pushs))
	return pushs, err
}

func GetPushsUserIds(userIds []string) ([]string, error) {
	var pushs []string
	var err = PushTokenTable.Find(bson.M{
		"user_id": bson.M{"$in": userIds},
	}).Distinct("push_token", &pushs)
	fmt.Println(len(pushs))
	return pushs, err
}

func GetFromToken(token string) *PushToken {
	if len(token) < 8 {
		panic(rest.Unauthorized("missing or invalid access token"))
	}
	var push *PushToken
	var err = PushTokenTable.FindByID(token, *push)
	rest.AssertNil(err)
	if push == nil || push.IsRevoke {
		rest.AssertNil(rest.Unauthorized("Hết thời gian truy cập"))
	}
	return push
}
