package mrw_booking

import (
	// "cetm_booking/o/ticket_onl"
	//"bytes"
	// "cetm_booking/x/fcm"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	//"encoding/json"
	"github.com/gin-gonic/gin"
	//"net/http"
)

type BookingServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewBookingServer(parent *gin.RouterGroup, name string) {
	var s = BookingServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/search", s.handlerSearch)
	s.GET("/services", s.handleService)
	s.POST("/ticket/create", s.handlerCreateTicket)
	s.POST("/ticket/cus_update", s.handlerUpdateTicketCus)
	s.POST("/ticket/cetm_update", s.handlerUpdateTicketCetm)
	s.POST("/ticket/canceled", s.handlerCancelTicket)
}

func (s *BookingServer) handlerSearch(ctx *gin.Context) {
	var body = []AddressBank{}
	rest.AssertNil(ctx.BindJSON(&body))
	var urlStr = "http://123.31.12.147:8888/room/booking/search_banks"
	s.SendData(ctx, web.ResParamArrUrlClient(urlStr, &body))
}

func (s *BookingServer) handleService(ctx *gin.Context) {
	var url = "http://mqserver:3000/room/booking/search_services"
	s.SendData(ctx, web.ResUrlClientGet(url))
}

type AddressBank struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
}

type InfoBankNow struct {
	AddressBank
	CountPeople int `json:"count_people"`
}
