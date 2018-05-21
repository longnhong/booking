package api

import (
	"cetm_booking/api/admin"
	"cetm_booking/api/auth"
	"cetm_booking/api/mrw_booking"
	"cetm_booking/system"
	"github.com/gin-gonic/gin"
)

func InitApi(root *gin.RouterGroup, tkWorker *system.TicketWorker) {
	mrw_booking.NewBookingServer(root, "booking", tkWorker)
	admin.NewAdminServer(root, "admin")
	auth.NewAuthenServer(root, "auth")
}
