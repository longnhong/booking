package mrw_booking

import (
	// "cetm_booking/middleware"
	// "cetm_booking/o/ticket_onl"
	//"bytes"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	//"encoding/json"
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
	s.POST("/search", s.handlerSearch)
	s.GET("/services", s.handleService)
}

func (s *BookingServer) handlerSearch(ctx *gin.Context) {
	var body = []AddressBank{}
	ctx.BindJSON(&body)
	var urlStr = "http://mqserver:3000/room/booking/search_banks"
	// /var banks = []InfoBankNow{}
	// result, err := json.Marshal(body)
	// req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(result))
	// req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// json.NewDecoder(resp.Body).Decode(&banks)

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
