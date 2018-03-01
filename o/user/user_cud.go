package user

import (
	"cetm_booking/o/push_token"
	"cetm_booking/x/rest"
)

func LoginApp(lg *LoginUser) (*User, string) {
	var err, res = GetUserByLogin(lg.Username, string(lg.Password))
	rest.AssertNil(err)
	return res, CreatePushToken(int(ROLE_USER), res.ID, lg.DeviceId, lg.PushToken).ID
}

func CreatePushToken(role int, userId string, deviceID string, pushToken string) *push_token.PushToken {
	var psh = push_token.PushToken{
		Role:      role,
		UserId:    userId,
		DeviceId:  deviceID,
		PushToken: pushToken,
	}
	return psh.CratePushToken()
}
