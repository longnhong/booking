package mrw_booking

import (
	"cetm_booking/api/mrw_booking/notify"
	"cetm_booking/api/mrw_booking/ticket"
	"cetm_booking/common"
	ctrl "cetm_booking/ctrl_to_cetm"
	user "cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/system"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type bookingServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewBookingServer(parent *gin.RouterGroup, name string, tkWorker *system.TicketWorker) {
	var s = bookingServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/search_branchs", s.handlerSearchs)
	s.GET("/search_services", s.handleService)
	ticket.NewTicketServer(s.RouterGroup, "ticket", tkWorker)
	notify.NewNotifyServer(s.RouterGroup, "notify")
}

func (s *bookingServer) handlerSearchs(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body = struct {
		KmScan       float64               `json:"km_scan"`
		ServiceID    string                `json:"service_id"`
		AddressBanks []*ctrl.AddressBank   `json:"address_banks"`
		TimeStart    int64                 `json:"time_start"`
		TimeEnd      int64                 `json:"time_end"`
		TypeSearch   ticket_onl.TypeTicket `json:"type_ticket"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	fmt.Println("CETM: " + common.ConfigSystemBooking.LinkCetm)
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_banks"
	var kmScan = common.ConfigSystemBooking.KmSearch
	body.KmScan = kmScan
	var res *ctrl.Data
	rest.AssertNil(web.ResParamArrUrlClient(urlStr, &body, &res))
	if res.Status == "error" {
		rest.AssertNil(errors.New("Có lỗi xảy ra!"))
	}
	var result = res.Data
	if body.TypeSearch == ticket_onl.TypeSchedule {
		branchIds := make([]string, len(result))
		for i, item := range result {
			branchIds[i] = item.BranchID
		}
		var tks, err = ticket_onl.SearchTicket(branchIds, body.TimeStart, body.TimeEnd)
		if err != nil {
			rest.AssertNil(errors.New("Có lỗi xảy ra!"))
		}
		for _, item := range result {
			var count = 0
			for _, tk := range tks {
				if tk.IdBranch == item.BranchID {
					count = tk.CountPeople
					break
				}
			}
			item.CountPeople = count
		}
	}
	s.SendData(ctx, result)
}

func (s *bookingServer) handleService(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	fmt.Println("CETM: " + common.ConfigSystemBooking.LinkCetm)
	var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_services"
	var res *ctrl.DataServices
	rest.AssertNil(web.ResUrlClientGet(url, &res))
	s.SendData(ctx, &res.Data)
}
