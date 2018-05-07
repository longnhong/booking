package push_token

import (
	"cetm_booking/x/mlog"
	"cetm_booking/x/rest"
	"cetm_booking/x/rest/validator"

	"gopkg.in/mgo.v2/bson"
)

var log = mlog.NewTagLog("push_token")

func (tok *PushToken) create() {
	rest.AssertNil(validator.Validate(tok))
	tok.IsRevoke = false
	// tok.BeforeCreate("k", 80)
}

func (tok *PushToken) delete() {
	tok.BeforeDelete()
}

func (tok *PushToken) update() {
	tok.BeforeUpdate()
	tok.IsRevoke = true
}

func checkTokenRevoke(token string) (*PushToken, error) {
	var tok PushToken
	return &tok, PushTokenTable.FindOne(bson.M{
		"id":        token,
		"is_revoke": true,
	}, &tok)
}

func CheckTokenByUserId(userId string) (*PushToken, error) {
	var tok PushToken
	return &tok, PushTokenTable.FindOne(bson.M{
		"user_id":   userId,
		"is_revoke": false,
	}, &tok)
}
