package auth

import (
	"cetm_booking/o/auth/user"
	"cetm_booking/o/push_token"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"net/http"
)

func GetUserFromToken(r *http.Request) (usr *user.User, push *push_token.PushToken) {
	var token = web.GetToken(r)
	push = push_token.GetFromToken(token)
	var err error
	usr, err = user.GetByID(push.UserId)
	rest.AssertNil(err)
	return
}

type UserAuth user.User

func GetUserByID(userID string) *user.User {
	var usr, err = user.GetByID(userID)
	rest.AssertNil(err)
	return usr
}

func GetFromToken(r *http.Request) {
	var token = web.GetToken(r)
	push_token.GetFromToken(token)
}

type LoginUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	DeviceId  string `json:"device_id"`
	PushToken string `json:"push_token"`
}

type UserPush struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}

type User struct {
	*user.User
}

func LoginApp(lg *LoginUser) *UserPush {
	var res, err = user.GetUserByLogin(lg.Username, lg.Password)
	rest.AssertNil(err)
	var userPush = UserPush{
		Token: CreatePushToken(int(user.ROLE_USER), res.ID, lg.DeviceId, lg.PushToken),
		User:  res,
	}
	return &userPush
}

func CreatePushToken(role int, userId string, deviceID string, pushToken string) string {
	var psh = push_token.PushToken{
		Role:      role,
		UserId:    userId,
		DeviceId:  deviceID,
		PushToken: pushToken,
	}
	return psh.CratePushToken().ID
}

func CreateUser(u *User) error {
	return u.CreateUser()
}
