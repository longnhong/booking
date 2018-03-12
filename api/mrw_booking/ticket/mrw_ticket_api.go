package ticket

import (
	"cetm_booking/common"
	user "cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
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

func NewTicketServer(parent *gin.RouterGroup, name string) {
	var s = TicketServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/create", s.handlerCreateTicket)
	s.POST("/get_ticket_day", s.handlerGetTicketDay)
	s.POST("/cus_update", s.handlerUpdateTicketCus)
	s.POST("/cetm_update", s.handlerUpdateTicketCetm)
	s.POST("/canceled", s.handlerCancelTicket)
	s.POST("/check_code", s.handlerCheckCode)
	s.GET("/branch_tickets", s.handlerGetTicketDayInBranch)
	s.POST("/check_location", s.handlerLoction)
}

func (s *TicketServer) handlerCreateTicket(ctx *gin.Context) {
	var userTK = user.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.CustomerID = userTK.ID
	var ticket = body.CrateTicketBooking()
	s.SendData(ctx, ticket)
}

func (s *TicketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk = body.UpdateTicketBookingByCustomer()
	s.SendData(ctx, tk)
}

func (s *TicketServer) handlerGetTicketDayInBranch(ctx *gin.Context) {
	var request = ctx.Request
	user.GetFromToken(request)
	var branchID = request.URL.Query().Get("branch_id")
	var reslt, err = ticket_onl.GetTicketDayInBranch(branchID)
	var result = make([]int64, len(reslt))
	for i, item := range reslt {
		result[i] = item.TimeGoBank
	}
	rest.AssertNil(err)
	var data = SearchBank(branchID)
	var res = map[string]interface{}{
		"bank":         data,
		"time_tickets": result,
	}
	s.SendData(ctx, res)
}

func (s *TicketServer) handlerUpdateTicketCetm(ctx *gin.Context) {
	var body = ticket_onl.UpdateCetm{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.UpdateTicketBookingByCetm()
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerCancelTicket(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body = struct {
		BTicketID string `json:"bticket_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	fmt.Println("TICKET :" + body.BTicketID)
	ticket_onl.CancleTicket(body.BTicketID)
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerGetTicketDay(ctx *gin.Context) {
	var usrTk = user.GetUserFromToken(ctx.Request)
	var bTks, err = ticket_onl.CheckTicketByDay(usrTk.ID)
	rest.AssertNil(err)
	s.SendData(ctx, bTks)
}

func (s *TicketServer) handlerCheckCode(ctx *gin.Context) {
	//var usrTk = user.GetUserFromToken(ctx.Request)
	var body = struct {
		CustomerCode string `json:"customer_code"`
		BranchId     string `json:"branch_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk, err = ticket_onl.CheckCustomerCode(body.CustomerCode, body.BranchId)
	rest.AssertNil(err)
	if tk == nil {
		rest.AssertNil(errors.New("Code sai"))
	}
	rest.AssertNil(tk.UpdateTimeCheckIn())
	s.SendData(ctx, tk)
}

func (s *TicketServer) handlerLoction(ctx *gin.Context) {
	var userTK = user.GetUserFromToken(ctx.Request)
	var body = struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var bTk, err = ticket_onl.CheckTicketByDay(userTK.ID)
	rest.AssertNil(err)
	var data = SearchBank(bTk.BranchID)
	if data == nil {
		rest.AssertNil(errors.New("Thử lại! Bạn đang không trong khu vực ngân hàng!"))
	}
	var distance = utility.Haversine(data.Lat, data.Lng, body.Lat, body.Lng)
	if distance < 0.02 {
		var noti = fcm.FmcMessage{
			Title: "Có việc mới!",
			Body:  "Bạn đang trong phạm vi ngân hàng!"}
		fcm.FcmCustomer.SendToOne("fmB1I_-GMqY:APA91bFVDqGQNKHcnYca6zzQ0ZG0kOyu92bloOcynHU4izvFFXvVbRIWglI2fVq4zp0XDILv282sQcTcX72lElG2VsmbfTzENj5rE_3R7pVCae8J57xaevCbXKGrZgzqwJnirembyUlM", noti)
		s.SendData(ctx, data)
	} else {
		rest.AssertNil(errors.New("Thử lại! Bạn đang không trong khu vực ngân hàng!"))
	}

}

func SearchBank(branchID string) *InfoBank {
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_bank?branch_id=" + branchID
	var data *DataBank
	rest.AssertNil(web.ResUrlClientGet(urlStr, &data))
	if data.Status == "error" {
		rest.AssertNil(errors.New("Không tìm thấy Branch này!"))
	}
	return data.Data
}
