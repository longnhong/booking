package middleware

import (
	//"cetm_booking/o/push_token"
	"cetm_booking/x/rest"
	// "fmt"
	// "g/x/web"
	"cetm_booking/x/mlog"
	"github.com/gin-gonic/gin"
)

var logMiddle = mlog.NewTagLog("middle")

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logMiddle.Error(err)
				var errResponse = map[string]interface{}{
					"error":  err.(error).Error(),
					"status": "error",
				}
				if httpError, ok := err.(rest.IHttpError); ok {
					errResponse["code"] = httpError.StatusCode()
					c.JSON(httpError.StatusCode(), errResponse)
				} else {
					errResponse["code"] = 500
					c.JSON(500, errResponse)
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func AddHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Range, Content-Disposition, Authorization",
		)
		c.Writer.Header().Set(
			"Access-Control-Allow-Credentials",
			"true",
		)
		//remember
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(200)
			return
		}
		c.Next()
	}
}

// func AuthenticateCustomer() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticateApp(c, 2)
// 	}
// }

// func MustAuthenticateApp(ctx *gin.Context, role int) {
// 	var errResponse = map[string]interface{}{
// 		"status": "error",
// 	}
// 	var token = web.GetToken(ctx.Request)
// 	var auth, err = push_token.GetByID(token)
// 	if err != nil {
// 		errResponse["error"] = "access token not found"
// 		ctx.JSON(401, errResponse)
// 	} else {
// 		if int(auth.Role) != role {
// 			errResponse["error"] = fmt.Sprintf("Unauthorize! you must be %s to access", role)
// 			ctx.JSON(401, errResponse)
// 		} else {
// 			ctx.Next()
// 		}
// 	}
// 	ctx.Abort()
// }
