package api

import (
	"cetm_booking/api/mrw_booking"
	"cetm_booking/api/user"
	"github.com/gin-gonic/gin"
)

func InitApi(root *gin.RouterGroup) {
	user.NewUserServer(root, "user")
	mrw_booking.NewBookingServer(root, "booking")
}
