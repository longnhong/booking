package mrw_booking

import (
	// "cetm_booking/o/ticket_onl"
	//"bytes"
	// "cetm_booking/x/fcm"
	ticket "cetm_booking/api/mrw_booking/ticket"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"errors"
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
	s.POST("/search_branchs", s.handlerSearchs)
	s.POST("/search_services", s.handleService)
	ticket.NewTicketServer(s.RouterGroup, "ticket")
}

func (s *BookingServer) handlerSearchs(ctx *gin.Context) {
	var body = []AddressBank{}
	rest.AssertNil(ctx.BindJSON(&body))
	var urlStr = "http://123.31.12.147:8888/room/booking/search_banks"
	var res *Data
	rest.AssertNil(web.ResParamArrUrlClient(urlStr, &body, &res))
	if res.Status == "error" {
		rest.AssertNil(errors.New("Có lỗi xảy ra!"))
	}
	s.SendData(ctx, res.Data)
}

func (s *BookingServer) handleService(ctx *gin.Context) {
	var url = "http://mqserver:3000/room/booking/search_services"
	var res interface{}
	rest.AssertNil(web.ResUrlClientGet(url, &res))
	s.SendData(ctx, &res)
}

type Data struct {
	Data   []*InfoBank `json:"data"`
	Status string      `json:"status"`
}

type DataBank struct {
	Data   InfoBank `json:"data"`
	Status string   `json:"status"`
}

type InfoBank struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Address     string  `json:"address"`
	BranchID    string  `json:branch_id`
	CountPeople int     `json:"count_people"`
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
