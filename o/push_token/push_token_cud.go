package push_token

import (
	"cetm_booking/common"
	"cetm_booking/x/rest"
)

func (tok *PushToken) CratePushToken() *PushToken {
	var res, err = CheckTokenByUserId(tok.UserId)
	if res != nil && err == nil {
		res.update()
		rest.AssertNil(PushTokenTable.UpdateId(res.ID, res))
	}
	tok.create()
	rest.AssertNil(PushTokenTable.Create(tok))
	return tok
}

func UpdatePushToken(tokenStr string) error {
	var res, err = checkTokenRevoke(tokenStr)
	if err != nil && err.Error() != common.NOT_EXIST {
		rest.AssertNil(err)
	}
	res.update()
	return PushTokenTable.UpdateId(res.ID, res)
}
