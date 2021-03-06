package ticket

import (
	"cetm_booking/common"
	ctrl "cetm_booking/ctrl_to_cetm"
	"cetm_booking/o/auth"
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerCreateTkCetm(ctx *gin.Context) {
	fmt.Println("CREATE TICKET CETM")
	var userTK, _ = auth.GetUserFromToken(ctx.Request)
	var body *common.Location
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var ticket = s.actionChange("", userTK.ID, ticket_onl.BookingStateCreateCetm, extra)
	ctrl.CreateNumCetm(userTK, ticket)
	s.SendData(ctx, ticket)
}
