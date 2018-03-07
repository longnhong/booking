package push_token

import (
	"cetm_booking/common"
	"cetm_booking/x/rest"
	"gopkg.in/mgo.v2/bson"
)

func (tok *PushToken) CratePushToken() *PushToken {
	var err, res = CheckTokenByUserId(tok.UserId)
	if res != nil && err == nil {
		res.update()
		rest.AssertNil(PushTokenTable.UpdateId(res.ID, res))
	}
	tok.create()
	rest.AssertNil(PushTokenTable.Create(tok))
	return tok
}

func UpdatePushToken(tokenStr string) error {
	var err, res = CheckTokenRevoke(tokenStr)
	if err != nil && err.Error() != common.NOT_EXIST {
		rest.AssertNil(err)
	}
	res.update()
	return PushTokenTable.UpdateId(res.ID, res)
}

func DeleteTokenByUser(userID string) error {
	var _, err = PushTokenTable.RemoveAll(bson.M{"user_id": userID})
	return err
}
