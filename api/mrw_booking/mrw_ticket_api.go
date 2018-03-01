package mrw_booking

import (
	"cetm_booking/o/ticket_onl"
	"cetm_booking/x/rest"
	"github.com/gin-gonic/gin"
)

func (s *BookingServer) handlerCreateTicket(ctx *gin.Context) {
	var body = ticket_onl.TicketBookingCreate{}
	rest.AssertNil(ctx.BindJSON(&body))
	var ticket = body.CrateTicketBooking()
	s.SendData(ctx, ticket)
}

func (s *BookingServer) handlerUpdateTicketCus(ctx *gin.Context) {
	var body = ticket_onl.TicketUpdate{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.UpdateTicketBookingByCustomer()
	s.SendData(ctx, nil)
}

func (s *BookingServer) handlerUpdateTicketCetm(ctx *gin.Context) {
	var body = ticket_onl.UpdateCetm{}
	rest.AssertNil(ctx.BindJSON(&body))
	body.UpdateTicketBookingByCetm()
	s.SendData(ctx, nil)
}

func (s *BookingServer) handlerCancelTicket(ctx *gin.Context) {
	var body struct {
		TicketID string `ticket_id`
	}
	rest.AssertNil(ctx.BindJSON(&body))
	ticket_onl.CancleTicket(body.TicketID)
	s.SendData(ctx, nil)
}
