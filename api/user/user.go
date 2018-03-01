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
	s.GET("/super_admin", s.handleExistSuperAdmin)
	s.POST("/create", s.handleCreate)
}

func (s *UserServer) handleCreate(ctx *gin.Context) {
	s.createUser(ctx)
}
func (s *UserServer) createUser(ctx *gin.Context) {
	var u *user.User
	ctx.BindJSON(&u)
	//rest.AssertNil(u.Create())
	s.SendData(ctx, u)
}
func (s *UserServer) handleExistSuperAdmin(ctx *gin.Context) {
	//s.SendData(ctx, user.GetSuperUser())
}
