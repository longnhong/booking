package rest

import (
	"cetm_booking/x/rest/validator"

	"github.com/gin-gonic/gin"
)

const STATUS_OK = 200

type JsonRender struct {
}

func (r *JsonRender) SendData(ctx *gin.Context, data interface{}) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   data,
		"status": "ok",
		"code":   200,
	})
}
func (r *JsonRender) DecodeBody(ctx *gin.Context, data interface{}) {
	AssertNil(BadRequest(ctx.BindJSON(&data).Error()))
	AssertNil(BadRequest(validator.Validate(data).Error()))

}
func (r *JsonRender) Success(ctx *gin.Context) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   nil,
		"status": "ok",
		"code":   200,
	})
}
