package ticket

import (
	"cetm_booking/common"
	user "cetm_booking/o/auth"
	oUser "cetm_booking/o/auth/user"
	"cetm_booking/o/push_token"
	"cetm_booking/o/rate"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/fcm"
	"cetm_booking/x/math"
	"cetm_booking/x/mlog"
	"cetm_booking/x/rest"
	"cetm_booking/x/utility"
	"cetm_booking/x/web"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
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
	s.GET("/send_push", s.handlerSendPush)
	s.GET("/ticket_priority", s.handlerPrioritys)
}

func (s *TicketServer) handlerSendPush(ctx *gin.Context) {
	var timeNow = time.Now()
	var tickets, _ = ticket_onl.GetAllTicketDay()

	for _, item := range tickets {
		if math.CompareDayTime(timeNow, item.TimeGoBank) == 0 {
			var dateTicket = time.Unix(item.TimeGoBank, 0)
			var pushTokens, _ = push_token.GetPushsUserId(item.CustomerID)
			var noti = fcm.FmcMessage{
				Title: "Gần đến giờ đặt vé!",
				Body:  "Thời gian: " + strconv.Itoa(dateTicket.Hour()) + ":" + strconv.Itoa(dateTicket.Minute()) + ".\nTrân trọng kính mời quý khách!",
			}
			fcm.FcmCustomer.SendToMany(pushTokens, noti)
		}
	}
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

func (s *TicketServer) handlerCreateTicket(ctx *gin.Context) {
	var userTK, push = user.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.CustomerID = userTK.ID
	var ticket = body.CrateTicketBooking()
	var countPP int
	SendFee(push.PushToken, ticket)
	if ticket.TypeTicket == ticket_onl.TYPE_NOW {
		countPP = CreateNumCetm(userTK, ticket)
	} else {
		var tks, _ = ticket_onl.GetAllTicketByTimeSearch(body.TimeGoBank)
		if tks != nil {
			countPP = len(tks)
		} else {
			countPP = 0
		}
	}
	var res = resData{
		TicketBooking: ticket,
		CountPeople:   countPP,
	}
	s.SendData(ctx, res)
}

type resData struct {
	*ticket_onl.TicketBooking
	CountPeople int `json:"count_people"`
}

func (s *TicketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk, err = ticket_onl.GetByID(body.BTicketID)
	rest.AssertNil(err)
	var usr *oUser.User
	if body.TypeTicket == ticket_onl.TYPE_NOW && tk.TypeTicket == ticket_onl.TYPE_SCHEDUCE {
		usr, _ = user.GetUserFromToken(ctx.Request)
	} else {
		user.GetFromToken(ctx.Request)
	}
	var tkUp = body.UpdateTicketBookingByCustomer()
	if usr != nil {
		CreateNumCetm(usr, tkUp)
	}
	s.SendData(ctx, tkUp)
}

type resTime struct {
	TimeGoBank int64                 `json:"time_go_bank"`
	TypeTicket ticket_onl.TypeTicket `json:"type_ticket"`
	ServiceID  string                `json:"service_id"`
}

func (s *TicketServer) handlerGetTicketDayInBranch(ctx *gin.Context) {
	var request = ctx.Request
	user.GetFromToken(request)
	var branchID = request.URL.Query().Get("branch_id")
	var reslt, err = ticket_onl.GetTicketDayInBranch(branchID)

	var result = make([]resTime, len(reslt))
	for i, item := range reslt {
		var res = resTime{
			TimeGoBank: item.TimeGoBank,
			TypeTicket: item.TypeTicket,
			ServiceID:  item.ServiceID,
		}
		result[i] = res
	}
	rest.AssertNil(err)
	var data = SearchBank(branchID)
	var res = map[string]interface{}{
		"bank":    data,
		"tickets": result,
	}
	s.SendData(ctx, res)
}

func (s *TicketServer) handlerGetTicketsDay(ctx *gin.Context) {
	var request = ctx.Request
	var branchID = request.URL.Query().Get("branch_id")
	var reslt, err = ticket_onl.GetTicketDayInBranch(branchID)
	rest.AssertNil(err)
	s.SendData(ctx, reslt)
}

func (s *TicketServer) handlerUpdateTicketCetm(ctx *gin.Context) {
	var body = ticket_onl.UpdateCetm{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk, err = ticket_onl.GetByID(body.BTicketID)
	rest.AssertNil(err)
	if tk.Status == ticket_onl.BOOKING_STATE_CANCELLED || ticket_onl.BOOKING_STATE_FINISHED == tk.Status {
		rest.AssertNil(errors.New("Vé đã được phản hồi"))
	}
	fmt.Printf("Thời gian đợi", body.WaitingTime)
	body.UpdateTicketBookingByCetm()
	pDevices, err := push_token.GetPushsUserId(tk.CustomerID)
	rest.AssertNil(err)
	fmt.Printf("SỐ DEVICE", pDevices)
	SendFeedback(pDevices, tk)
	s.SendData(ctx, nil)
}

func (s *TicketServer) handlerDeleteTicket(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body = struct {
		BTicketID string `json:"bticket_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	fmt.Println("TICKET :" + body.BTicketID)
	var tk = ticket_onl.MarkDeleteTicket(body.BTicketID)
	if tk.TypeTicket == ticket_onl.TYPE_NOW {
		var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/cancel_bticket"
		fmt.Println(url)
		var data = Status{}
		rest.AssertNil(web.ResParamArrUrlClient(url, tk, &data))
		if data.Status == "error" {
			logApi.Errorf("handlerDeleteTicket ", "Error CeTM /room/booking/cancel_bticket")
		}
	}
	s.SendData(ctx, nil)
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

func (s *TicketServer) handlerCheckCode(ctx *gin.Context) {
	//var usrTk = user.GetUserFromToken(ctx.Request)
	var body = struct {
		CustomerCode string `json:"customer_code"`
		BranchId     string `json:"branch_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var ticket, err = ticket_onl.CheckCustomerCode(body.CustomerCode, body.BranchId)
	rest.AssertNil(err)
	if ticket == nil {
		rest.AssertNil(errors.New("Code sai"))
	}
	rest.AssertNil(ticket.UpdateTimeCheckIn())
	var userTK = user.GetUserByID(ticket.CustomerID)
	var data = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	s.SendData(ctx, data)
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
	var data = SearchBank(bTk.BranchID)
	if data == nil {
		rest.AssertNil(errors.New("Thử lại! Bạn đang không trong khu vực ngân hàng!"))
	}
	var distance = utility.Haversine(data.Lat, data.Lng, body.Lat, body.Lng)
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

func SearchBank(branchID string) *InfoBank {
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_bank?branch_id=" + branchID
	var data *DataBank
	rest.AssertNil(web.ResUrlClientGet(urlStr, &data))
	if data.Status.Status == "error" {
		rest.AssertNil(errors.New("Không tìm thấy Branch này!"))
	}
	return data.Data
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
