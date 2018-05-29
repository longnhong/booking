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
	s.GET("/booking_detail", s.ReportDetail)
}

func (s *adminServer) ReportPerformance(ctx *gin.Context) {
	var request = ctx.Request
	var branchIds = web.GetArrString("branch_id", ",", request.URL.Query())
	var timeStart = web.MustGetInt64("start", request.URL.Query())
	var timeEnd = web.MustGetInt64("end", request.URL.Query())
	var skip = web.MustGetInt64("skip", request.URL.Query())
	var limit = web.MustGetInt64("limit", request.URL.Query())
	var tks, err = ticket_onl.GetTicketReport(branchIds, timeStart, timeEnd, skip, limit)
	rest.AssertNil(err)
	var countAll, err1 = ticket_onl.GetTicketReportByTime(branchIds, timeStart, timeEnd)
	rest.AssertNil(err1)
	var data = map[string]interface{}{
		"tickets": tks,
		"total":   countAll,
	}
	s.SendData(ctx, data)
}

func (s *adminServer) ReportDetail(ctx *gin.Context) {
	var request = ctx.Request
	var branchIds = web.GetArrString("branch_id", ",", request.URL.Query())
	var timeStart = web.MustGetInt64("start", request.URL.Query())
	var timeEnd = web.MustGetInt64("end", request.URL.Query())
	// var skip = web.MustGetInt64("skip", request.URL.Query())
	// var limit = web.MustGetInt64("limit", request.URL.Query())
	var tks, err = ticket_onl.GetDetailReport(branchIds, timeStart, timeEnd)
	rest.AssertNil(err)
	s.SendData(ctx, tks)
}
