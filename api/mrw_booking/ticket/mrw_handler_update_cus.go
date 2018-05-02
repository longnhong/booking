package ticket

import (
	"cetm_booking/o/auth"
	"cetm_booking/o/auth/user"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"github.com/gin-gonic/gin"
)

func (s *TicketServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))

	var tk, err = ticket_onl.GetByID(body.BTicketID)
	rest.AssertNil(err)
	var usr *user.User
	if body.TypeTicket == ticket_onl.TYPE_NOW && tk.TypeTicket == ticket_onl.TYPE_SCHEDULE {
		usr, _ = auth.GetUserFromToken(ctx.Request)
	} else {
		auth.GetFromToken(ctx.Request)
	}
	var tkUp = body.UpdateTicketBookingByCustomer()
	if usr != nil {
		CreateNumCetm(usr, tkUp)
	}

	s.SendData(ctx, tk)
}
