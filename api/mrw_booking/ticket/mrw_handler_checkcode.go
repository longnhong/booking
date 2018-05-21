package ticket

import (
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerCheckCode(ctx *gin.Context) {
	//var usrTk = user.GetUserFromToken(ctx.Request)
	fmt.Println("CHECK CODE")
	var body = struct {
		CustomerCode string `json:"customer_code"`
		BranchID     string `json:"branch_id"`
	}{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var ticket = s.actionChange("", "", ticket_onl.BOOKING_STATE_CHECK_CODE, extra)
	if ticket == nil {
		rest.AssertNil(errors.New("Code sai"))
	}
	var userTK = auth.GetUserByID(ticket.CustomerID)
	var data = ctrl.DataTicketSendCetm{
		TicketBooking: ticket,
		Customer:      userTK,
	}
	s.SendData(ctx, data)
}
