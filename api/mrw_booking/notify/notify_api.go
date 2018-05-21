package notify

import (
	"cetm_booking/o/auth"
	"cetm_booking/o/notify"
	"cetm_booking/x/mlog"
	"cetm_booking/x/rest"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NotifyServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

var logApi = mlog.NewTagLog("ticket_API")

func NewNotifyServer(parent *gin.RouterGroup, name string) {
	var s = NotifyServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/read", s.handlerReadNotify)
	s.GET("/get", s.handlerGetNotify)
	s.GET("/get_readnt", s.handlerGetNoRead)
	s.POST("/update_reads", s.handlerUpdateReadeds)
}

func (s *NotifyServer) handlerGetNoRead(ctx *gin.Context) {
	var psh = auth.GetFromToken(ctx.Request)
	var count, err = notify.CountNotifyNoRead(psh.UserId)
	rest.AssertNil(err)
	s.SendData(ctx, count)
}

func (s *NotifyServer) handlerUpdateReadeds(ctx *gin.Context) {
	var psh = auth.GetFromToken(ctx.Request)
	var body = struct {
		Ids []string `json:"notify_ids"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var err = notify.UpdateNotifyReaded(psh.UserId, body.Ids)
	rest.AssertNil(err)
	s.SendData(ctx, nil)
}

func (s *NotifyServer) handlerGetNotify(ctx *gin.Context) {
	var cus, _ = auth.GetUserFromToken(ctx.Request)
	var query = ctx.Request.URL.Query()
	var skip, _ = strconv.ParseInt(query.Get("skip"), 10, 64)
	var limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	var noties = notify.GetAllNotifyByCus(cus.ID, int(skip), int(limit))
	s.SendData(ctx, noties)
}

func (s *NotifyServer) handlerReadNotify(ctx *gin.Context) {
	auth.GetFromToken(ctx.Request)
	var body = struct {
		IdNotify string `json:"id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	rest.AssertNil(notify.UpdateRead(body.IdNotify))
	s.SendData(ctx, nil)
}
