package push_token

import (
	"cetm_booking/o/push_token"
	"cetm_booking/x/rest"
	"github.com/gin-gonic/gin"
)

type PushServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewPushServer(parent *gin.RouterGroup, name string) {
	var s = PushServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/search_branchs", s.handleLogin)
	//ath.NewAuthenServer(s.RouterGroup, "auth")
}

func (s *PushServer) handleLogin(ctx *gin.Context) {
	var token = "afasfasfasfasdf"
	var _, err = push_token.GetByID(token)
	s.SendData(ctx, err)
}
