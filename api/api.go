package api

import (
	"cetm_booking/api/auth"
	"cetm_booking/api/mrw_booking"
	"github.com/gin-gonic/gin"
)

func InitApi(root *gin.RouterGroup) {
	mrw_booking.NewBookingServer(root, "booking")
	auth.NewAuthenServer(root, "auth")
}
