package ticket

import (
	"cetm_booking/common"
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"cetm_booking/x/web"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerDeleteTicket(ctx *gin.Context) {
	fmt.Println("DELETE TICKET")
	auth.GetFromToken(ctx.Request)
	var body = struct {
		BTicketID string `json:"bticket_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk = actionChange(body.BTicketID, "", ticket_onl.BOOKING_STATE_DELETE, nil)
	if tk.TypeTicket == ticket_onl.TYPE_NOW {
		var url = common.ConfigSystemBooking.LinkCetm + "/room/booking/cancel_bticket"
		var data = struct {
			Status string `json:"status"`
		}{}
		rest.AssertNil(web.ResParamArrUrlClient(url, tk, &data))
		if data.Status == "error" {
			logAPI.Errorf("handlerDeleteTicket ", "Error CeTM /room/booking/cancel_bticket")
		}
	}
	s.SendData(ctx, nil)
}
