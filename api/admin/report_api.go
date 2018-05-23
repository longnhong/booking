package admin

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"github.com/gin-gonic/gin"
)

type adminServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAdminServer(parent *gin.RouterGroup, name string) {
	var s = adminServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("/appointment_performance", s.ReportPerformance)
	//s.GET("/excel_performance", s.ReportPerformance)
}

func (s *adminServer) ReportPerformance(ctx *gin.Context) {
	var request = ctx.Request
	var branchIds = web.GetArrString("branch_id", ",", request.URL.Query())
	var timeStart = web.MustGetInt64("start", request.URL.Query())
	var timeEnd = web.MustGetInt64("end", request.URL.Query())
	var tks, err = ticket_onl.GetTicketReport(branchIds, timeStart, timeEnd)
	rest.AssertNil(err)
	s.SendData(ctx, tks)
}
