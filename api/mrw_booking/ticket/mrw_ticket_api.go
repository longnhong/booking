package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	user "cetm_booking/o/auth"
	"cetm_booking/o/rate"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/system"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
	"cetm_booking/x/mlog"
	"cetm_booking/x/rest"
	"cetm_booking/x/ultility"
	"cetm_booking/x/web"
	"errors"
	"github.com/gin-gonic/gin"
)

type TicketServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

var logApi = mlog.NewTagLog("ticket_API")

func NewTicketServer(parent *gin.RouterGroup, name string) {
	var s = TicketServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/create", s.handlerCreateTicket)
	s.GET("/my_schedule", s.handlerMySchedule)
	s.GET("/mine_all", s.handlerGetTicketAll)
	s.POST("/cus_update", s.handlerUpdateTicketCus)
	s.POST("/cetm_update", s.handlerUpdateTicketCetm)
	s.POST("/delete", s.handlerDeleteTicket)
	s.POST("/check_code", s.handlerCheckCode)
	s.GET("/branch_tickets", s.handlerGetTicketDayInBranch)
	s.GET("/branch_cetm_tickets", s.handlerGetTicketsDay)
	s.POST("/check_location", s.handlerLoction)
	s.GET("/ticket_near", s.handlerTicketNear)
	s.POST("/rate", s.handlerRate)
	s.POST("/no_rate", s.handlerNoRate)
	s.GET("/get_ticket", s.handlerGetTicket)
	s.GET("/ticket_priority", s.handlerPrioritys)
	s.POST("/crate_tkcetm", s.handlerCreateTkCetm)
}

func (s *TicketServer) handlerGetTicket(ctx *gin.Context) {
	var request = ctx.Request
	user.GetFromToken(request)
	var btkID = request.URL.Query().Get("bticket_id")
	tk, err := system.GetTicketByID(btkID)
	rest.AssertNil(err)
	s.SendData(ctx, tk)
}

func (s *TicketServer) handlerPrioritys(ctx *gin.Context) {
	var userTK, _ = user.GetUserFromToken(ctx.Request)
	tks, err := ticket_onl.GetTicketByUserNeedFeedback(userTK.ID)
	rest.AssertNil(err)
	s.SendData(ctx, tks)
}

func (s *TicketServer) handlerGetTicketAll(ctx *gin.Context) {
	var userTK, _ = user.GetUserFromToken(ctx.Request)
	var tks, err = ticket_onl.GetAllTicketCus(userTK.ID)
	rest.AssertNil(err)
	s.SendData(ctx, tks)
}

func (s *TicketServer) handlerTicketNear(ctx *gin.Context) {
	var userTK, _ = user.GetUserFromToken(ctx.Request)
	var tks, err = ticket_onl.GetTicketNear(userTK.ID)
	rest.AssertNil(err)
	s.SendData(ctx, tks)
}

func (s *TicketServer) handlerGetTicketDayInBranch(ctx *gin.Context) {
	var request = ctx.Request
	user.GetFromToken(request)
	var branchID = request.URL.Query().Get("branch_id")
	var serviceID = request.URL.Query().Get("service_id")
	var timeStart = web.MustGetInt64("start", request.URL.Query())
	var timeEnd = web.MustGetInt64("end", request.URL.Query())
	var res, err = SetBankTickets(branchID, serviceID, timeStart, timeEnd)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

type bankTickets struct {
	Bank    *ctrl.InfoBank `json:"bank"`
	Tickets []resTime      `json:"tickets"`
}

func (s *TicketServer) handlerGetTicketsDay(ctx *gin.Context) {
	var request = ctx.Request
	var branchID = request.URL.Query().Get("branch_id")
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var reslt, err = ticket_onl.GetTicketDayInBranch(branchID, timeBeginDay, tiemEnOfday)
	rest.AssertNil(err)
	s.SendData(ctx, reslt)
}

func (s *TicketServer) handlerMySchedule(ctx *gin.Context) {
	var usrTk, _ = user.GetUserFromToken(ctx.Request)
	var bTks, err = ticket_onl.GetCustomerMySchedule(usrTk.ID)
	rest.AssertNil(err)
	s.SendData(ctx, bTks)
}

func (s *TicketServer) handlerRate(ctx *gin.Context) {
	var usrTk, _ = user.GetUserFromToken(ctx.Request)
	var body *rate.Rate
	rest.AssertNil(ctx.BindJSON(&body))
	var err = ticket_onl.UpdateRate(body.TicketIdBk, ticket_onl.TYPE_RATED)
	rest.AssertNil(err)
	body.CustomerId = usrTk.ID
	body.CrateRate()
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerNoRate(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body *rate.Rate
	var err = ticket_onl.UpdateRate(body.TicketIdBk, ticket_onl.TYPE_RATED)
	rest.AssertNil(err)
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerLoction(ctx *gin.Context) {
	var userTK, _ = user.GetUserFromToken(ctx.Request)
	var body = struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var bTk, err = ticket_onl.CheckTicketByDay(userTK.ID)
	rest.AssertNil(err)
	var data, err1 = ctrl.SearchBank(bTk.BranchID, "")
	rest.AssertNil(err1)
	if data == nil {
		rest.AssertNil(errors.New("Thử lại! Bạn đang không trong khu vực ngân hàng!"))
	}
	var distance = ultility.Haversine(data.Lat, data.Lng, body.Lat, body.Lng)
	if distance < 0.02 {
		var noti = fcm.FmcMessage{
			Title: "Thông báo!",
			Body:  "Bạn đang trong phạm vi ngân hàng!"}
		fcm.FcmCustomer.SendToOne("fmB1I_-GMqY:APA91bFVDqGQNKHcnYca6zzQ0ZG0kOyu92bloOcynHU4izvFFXvVbRIWglI2fVq4zp0XDILv282sQcTcX72lElG2VsmbfTzENj5rE_3R7pVCae8J57xaevCbXKGrZgzqwJnirembyUlM", noti)
		s.SendData(ctx, data)
	} else {
		rest.AssertNil(errors.New("Thử lại! Bạn đang không trong khu vực ngân hàng!"))
	}

}
