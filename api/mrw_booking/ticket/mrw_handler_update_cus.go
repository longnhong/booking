package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/auth"
	"cetm_booking/o/auth/user"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var ticket = s.actionChange(body.BTicketID, "", ticket_onl.BOOKING_CUSTOMER_UPDATE, extra)
	var usr *user.User
	if body.TypeTicket == ticket_onl.TYPE_NOW && ticket.TypeTicket == ticket_onl.TYPE_SCHEDULE {
		usr, _ = auth.GetUserFromToken(ctx.Request)
	} else {
		auth.GetFromToken(ctx.Request)
	}
	if usr != nil {
		ctrl.CreateNumCetm(usr, ticket)
	}
	s.SendData(ctx, ticket)
}
