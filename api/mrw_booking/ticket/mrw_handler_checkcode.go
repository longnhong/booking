package ticket

import (
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

func (s *TicketServer) handlerCheckCode(ctx *gin.Context) {
	//var usrTk = user.GetUserFromToken(ctx.Request)
	var body = struct {
		CustomerCode string `json:"customer_code"`
		BranchId     string `json:"branch_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var ticket = ActionChange("", "", ticket_onl.BOOKING_STATE_CHECK_CODE, extra)
	if ticket == nil {
		rest.AssertNil(errors.New("Code sai"))
	}
	var userTK = auth.GetUserByID(ticket.CustomerID)
	var data = DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	s.SendData(ctx, data)
}
