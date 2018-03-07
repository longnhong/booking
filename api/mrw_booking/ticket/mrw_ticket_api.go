package mrw_booking

import (
	"cetm_booking/o/ticket_onl"
	"fmt"
	//"cetm_booking/o/user"
	"cetm_booking/x/fcm"
	"cetm_booking/x/rest"
	"cetm_booking/x/utility"
	"cetm_booking/x/web"
	"errors"
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
	s.POST("/cus_update", s.handlerUpdateTicketCus)
	s.POST("/cetm_update", s.handlerUpdateTicketCetm)
	s.POST("/canceled", s.handlerCancelTicket)
	s.POST("/check_code", s.handlerCheckCode)
	s.POST("/check_location", s.handlerLoction)
}

func (s *TicketServer) handlerCreateTicket(ctx *gin.Context) {
	///	var userTK = user.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	//body.CustomerID = userTK.ID
	var ticket = body.CrateTicketBooking()
	s.SendData(ctx, ticket)
}

func (s *TicketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk = body.UpdateTicketBookingByCustomer()
	s.SendData(ctx, tk)
}

func (s *TicketServer) handlerUpdateTicketCetm(ctx *gin.Context) {
	var body = ticket_onl.UpdateCetm{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.UpdateTicketBookingByCetm()
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerCancelTicket(ctx *gin.Context) {
	//	user.GetFromToken(ctx.Request)
	var body = struct {
		BTicketID string `json:"bticket_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	fmt.Println("TICKET :" + body.BTicketID)
	ticket_onl.CancleTicket(body.BTicketID)
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerCheckCode(ctx *gin.Context) {
	//	var usrTk = user.GetUserFromToken(ctx.Request)
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
	s.SendData(ctx, tk)
}

func (s *TicketServer) handlerLoction(ctx *gin.Context) {
	//user.GetFromToken(ctx.Request)
	var body = struct {
		CustomerID string  `json:"customer_id"`
		Lat        float64 `json:"lat"`
		Lng        float64 `json:"lng"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var bTk, err = ticket_onl.CheckCustomerIdByDay(body.CustomerID)
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
	var urlStr = "http://123.31.12.147:8888/room/booking/search_bank?branch_id=" + branchID
	var data *DataBank
	rest.AssertNil(web.ResUrlClientGet(urlStr, &data))
	if data.Status == "error" {
		rest.AssertNil(errors.New("Không tìm thấy Branch này!"))
	}
	return data.Data
}

type DataBank struct {
	Data   *InfoBank `json:"data"`
	Status string    `json:"status"`
}

type InfoBank struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Address     string  `json:"address"`
	BranchID    string  `json:branch_id`
	CountPeople int     `json:"count_people"`
}
