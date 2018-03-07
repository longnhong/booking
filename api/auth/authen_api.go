package auth

import (
	oAuth "cetm_booking/o/auth"
	"cetm_booking/x/rest"
	"github.com/gin-gonic/gin"
)

type AuthenServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAuthenServer(parent *gin.RouterGroup, name string) {
	var s = AuthenServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/login", s.handleLogin)

}
func (s *AuthenServer) handleLogin(ctx *gin.Context) {
	var u *oAuth.LoginUser
	rest.AssertNil(ctx.BindJSON(&u))
	s.SendData(ctx, oAuth.LoginApp(u))
}
