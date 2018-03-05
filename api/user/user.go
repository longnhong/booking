package user

import (
	//"cetm_booking/middleware"
	"cetm_booking/o/user"
	"cetm_booking/x/rest"
	"github.com/gin-gonic/gin"
)

type UserServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewUserServer(parent *gin.RouterGroup, name string) {
	var s = UserServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/login", s.handleLogin)
}

func (s *UserServer) handleLogin(ctx *gin.Context) {
	var u *user.LoginUser
	rest.AssertNil(ctx.BindJSON(&u))
	s.SendData(ctx, u)
}
