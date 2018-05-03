package ticket

import (
	"cetm_booking/common"
	user "cetm_booking/o/auth"
	oUser "cetm_booking/o/auth/user"
	"cetm_booking/o/rate"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/system"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
	"cetm_booking/x/mlog"
	"cetm_booking/x/rest"
	"cetm_booking/x/utility"
	"cetm_booking/x/web"
	"errors"
	"fmt"
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

type resData struct {
	*ticket_onl.TicketBooking
	CountPeople int `json:"count_people"`
}

type resTime struct {
	ID         string                `json:"id"`
	TimeGoBank int64                 `json:"time_go_bank"`
	TypeTicket ticket_onl.TypeTicket `json:"type_ticket"`
	ServiceID  string                `json:"service_id"`
}

func (s *TicketServer) handlerGetTicketDayInBranch(ctx *gin.Context) {
	var request = ctx.Request
	user.GetFromToken(request)
	var branchID = request.URL.Query().Get("branch_id")
	var serviceID = request.URL.Query().Get("service_id")
	var timeStart = web.MustGetInt64("start", request.URL.Query())
	var timeEnd = web.MustGetInt64("end", request.URL.Query())
	var res = SetBankTickets(branchID, serviceID, timeStart, timeEnd)
	s.SendData(ctx, res)
}

type bankTickets struct {
	Bank    *InfoBank `json:"bank"`
	Tickets []resTime `json:"tickets"`
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
	var data = SearchBank(bTk.BranchID, "")
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

func SearchBank(branchID string, serviceID string) *InfoBank {
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_bank?branch_id=" + branchID + "&service_id=" + serviceID
	var data *DataBank
	rest.AssertNil(web.ResUrlClientGet(urlStr, &data))
	if data.Status.Status == "error" {
		rest.AssertNil(errors.New("Không tìm thấy Branch này!"))
	}
	return data.Data
}

func UpdateCounterTkCetm(userTK *oUser.User, ticket *ticket_onl.TicketBooking) (err error) {
	var dataTicketSend = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/update_bticket"
	var data = struct {
		Data   interface{} `json:"data"`
		Status string      `json:"status"`
	}{}
	rest.AssertNil(web.ResParamArrUrlClient(urlStr, dataTicketSend, &data))
	if data.Status != "error" {
		err = errors.New(data.Status)
	}
	return
}

func CreateTicket(tk *ticket_onl.TicketBooking) *InfoBank {
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/system_add_bkticket"
	var data *DataBank
	rest.AssertNil(web.ResUrlClientGet(urlStr, &data))
	if data.Status.Status == "error" {
		rest.AssertNil(errors.New("Không tìm thấy Branch này!"))
	}
	return data.Data
}

func CreateNumCetm(userTK *oUser.User, ticket *ticket_onl.TicketBooking) (countPP int) {
	var dataTicketSend = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/system_add_bkticket"
	var data = DataTicketBookNow{}
	rest.AssertNil(web.ResParamArrUrlClient(url, dataTicketSend, &data))
	ticket.UpdateByCnumCetm(data.Data.Cnum, data.Data.Id)
	fmt.Printf("Data tao ve ", data.Data)
	return data.Data.CountPeople
}
