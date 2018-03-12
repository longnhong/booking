package mrw_booking

import (
	ticket "cetm_booking/api/mrw_booking/ticket"
	"cetm_booking/common"
	user "cetm_booking/o/auth"
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

}

func (s *BookingServer) handlerSearchs(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	var body = []AddressBank{}
	rest.AssertNil(ctx.BindJSON(&body))
	fmt.Println("CETM: " + common.ConfigSystemBooking.LinkCetm)
	var urlStr = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_banks"
	var res *Data
	rest.AssertNil(web.ResParamArrUrlClient(urlStr, &body, &res))
	if res.Status == "error" {
		rest.AssertNil(errors.New("Có lỗi xảy ra!"))
	}
	s.SendData(ctx, res.Data)
}

func (s *BookingServer) handleService(ctx *gin.Context) {
	user.GetFromToken(ctx.Request)
	fmt.Println("CETM: " + common.ConfigSystemBooking.LinkCetm)
	var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/search_services"
	var res *DataServices
	rest.AssertNil(web.ResUrlClientGet(url, &res))
	s.SendData(ctx, &res.Data)
}
