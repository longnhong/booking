package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/auth"
	"cetm_booking/o/auth/user"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/math"
	"cetm_booking/x/rest"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var ticket = s.actionChange(body.BTicketID, "", ticket_onl.BookingCustomerUpdate, extra)
	var usr *user.User
	if body.TypeTicket == ticket_onl.TypeNow && ticket.TypeTicket == ticket_onl.TypeSchedule {
		usr, _ = auth.GetUserFromToken(ctx.Request)
		var timeNow = math.GetTimeNowVietNam()
		if math.CompareDayTime(timeNow, body.TimeGoBank) == 0 {
			ctrl.UpdateCounterTkCetm(usr, ticket)
		}
	} else {
		auth.GetFromToken(ctx.Request)
	}
	if usr != nil {
		ctrl.CreateNumCetm(usr, ticket, true)
	}
	s.SendData(ctx, ticket)
}
