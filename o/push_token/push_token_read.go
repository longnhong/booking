package push_token

import (
	"cetm_booking/x/rest"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*PushToken, error) {
	var auth *PushToken
	return auth, PushTokenTable.FindOne(bson.M{
		"_id":       id,
		"is_revoke": false,
	}, &auth)
}

type PushDevice struct {
	ID         string   `json:"device_id" bson:"_id"`
	PushTokens []string `json:"push_tokens" bson:"push_tokens"`
}

func GetPushsUserId(userId string) ([]string, error) {
	var pushs []*PushDevice
	var queryMatch = bson.M{"user_id": userId}
	var sortMatch = bson.M{"created_at": -1}
	var group = bson.M{"_id": "$device_id", "push_tokens": bson.M{"$push": "$push_token"}}
	var query = []bson.M{
		{"$match": queryMatch},
		{"$sort": sortMatch},
		{"$group": group},
	}
	err := PushTokenTable.Pipe(query).All(&pushs)
	if err != nil {
		return nil, err
	}
	var pushTokens []string
	for _, item := range pushs {
		if len(item.PushTokens) > 0 {
			pushTokens = append(pushTokens, item.PushTokens[0])
		}
	}
	return pushTokens, nil
}

func GetFromToken(token string) *PushToken {
	if len(token) < 8 {
		panic(rest.Unauthorized("missing or invalid access token"))
	}
	var push *PushToken
	var err = PushTokenTable.FindByID(token, &push)
	rest.IsErrorRecord(err)
	rest.AssertNil(err)
	if push == nil || push.IsRevoke {
		rest.AssertNil(rest.Unauthorized("Hết thời gian truy cập"))
	}
	return push
}
