package mrw_booking

import (
	"cetm_booking/api/mrw_booking/notify"
	ticket "cetm_booking/api/mrw_booking/ticket"
	"cetm_booking/common"
	user "cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type BookingServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewBookingServer(parent *gin.RouterGroup, name string) {
	var s = BookingServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/search_branchs", s.handlerSearchs)
	s.GET("/search_services", s.handleService)
	ticket.NewTicketServer(s.RouterGroup, "ticket")
	notify.NewNotifyServer(s.RouterGroup, "notify")
}

func (s *BookingServer) handlerSearchs(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body = struct {
		KmScan       float64               `json:"km_scan"`
		ServiceID    string                `json:"service_id"`
		AddressBanks []*AddressBank        `json:"address_banks"`
		TimeStart    int64                 `json:"time_start"`
		TimeEnd      int64                 `json:"time_end"`
		TypeSearch   ticket_onl.TypeTicket `json:"type_ticket"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	fmt.Println("CETM: " + common.ConfigSystemBooking.LinkCetm)
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_banks"
	var kmScan = common.ConfigSystemBooking.KmSearch
	body.KmScan = kmScan
	var res *Data
	rest.AssertNil(web.ResParamArrUrlClient(urlStr, &body, &res))
	if res.Status == "error" {
		rest.AssertNil(errors.New("Có lỗi xảy ra!"))
	}
	var result = res.Data
	if body.TypeSearch == ticket_onl.TYPE_SCHEDUCE {
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

func (s *BookingServer) handleService(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	fmt.Println("CETM: " + common.ConfigSystemBooking.LinkCetm)
	var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_services"
	var res *DataServices
	rest.AssertNil(web.ResUrlClientGet(url, &res))
	s.SendData(ctx, &res.Data)
}
