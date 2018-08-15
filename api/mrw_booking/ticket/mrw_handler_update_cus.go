package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var usr, _ = auth.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var ticket = s.actionChange(body.BTicketID, "", ticket_onl.BookingCustomerUpdate, extra)
	fmt.Printf("THỜI GIAN THAY DOI", body)
	var timeNow = math.GetTimeNowVietNam()
	fmt.Printf("THỜI GIAN THAY DOI", timeNow.Unix())
	if math.CompareDayTime(timeNow, body.TimeGoBank) == 0 {
		ctrl.UpdateCounterTkCetm(usr, ticket)
	}
	s.SendData(ctx, ticket)
}
