package ticket

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (s *ticketServer) handlerUpdateTicketCetm(ctx *gin.Context) {
	var body = ticket_onl.UpdateCetm{}
	rest.AssertNil(ctx.BindJSON(&body))
	var extra, _ = json.Marshal(body)
	var tk = actionChange(body.BTicketID, "", body.Status, extra)
	s.SendData(ctx, tk)
}
