package ticket

import (
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (s *TicketServer) handlerCreateTicket(ctx *gin.Context) {
	var userTK, push = auth.GetUserFromToken(ctx.Request)
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(map[string]interface{}{
		"ticket":     body,
		"push_token": push.PushToken})
	var ticket = ActionChange("", userTK.ID, ticket_onl.BOOKING_STATE_CREATED, extra)
	var countPP int
	if ticket.TypeTicket == ticket_onl.TYPE_NOW {
		countPP = CreateNumCetm(userTK, ticket)
	} else {
		var tks, _ = ticket_onl.GetAllTicketByTimeSearch(body.TimeGoBank, body.TypeTicket)
		if tks != nil {
			countPP = len(tks)
		} else {
			countPP = 0
		}
	}
	var res = resData{
		TicketBooking: ticket,
		CountPeople:   countPP,
	}
	s.SendData(ctx, res)
}
