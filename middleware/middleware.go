package middleware

import (
	//"cetm_booking/o/push_token"
	//	uAuth "cetm_booking/o/user/auth"
	"cetm_booking/x/rest"
	//"fmt"
	//"g/x/web"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				glog.Error(err)
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

// func Authenticate() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var token = web.GetToken(c.Request)
// 		var _, err = auth.GetByID(token)
// 		if err != nil {
// 			rest.AssertNil(rest.Unauthorized(err.Error()))
// 		}
// 	}
// }

// func AuthenticateAppRole(role uAuth.Role) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var token = web.GetToken(c.Request)
// 		var res, err = push_token.GetByID(token)
// 		if err != nil {
// 			rest.AssertNil(rest.Unauthorized(err.Error()))
// 		}
// 		if res.Role != int(role) {
// 			rest.AssertNil(rest.Unauthorized("Bạn không có quyền truy cập!"))
// 		}
// 	}
// }

// func MustBeAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticate(c, "admin")
// 	}
// }
// func MustBeSuperAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticate(c, "super-admin")
// 	}
// }
// func MustBeStaff() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticate(c, "staff")
// 	}
// }
// func MustBeOwner() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticate(c, "owner")
// 	}
// }
// func MustAuthenticate(ctx *gin.Context, role string) {
// 	var errResponse = map[string]interface{}{
// 		"status": "error",
// 	}
// 	var token = web.GetToken(ctx.Request)
// 	var auth, err = auth.GetByID(token)
// 	if err != nil {
// 		errResponse["error"] = "access token not found"
// 		ctx.JSON(401, errResponse)
// 	} else {
// 		if auth.Role != role {
// 			errResponse["error"] = fmt.Sprintf("Unauthorize! you must be %s to access", role)
// 			ctx.JSON(401, errResponse)
// 		} else {
// 			ctx.Next()
// 		}
// 	}
// 	ctx.Abort()
// }

// func MustBeRoleEmployee() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticateRole(c, uAuth.RoleEmployee)
// 	}
// }
// func MustBeRoleCustomer() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticateRole(c, uAuth.RoleCustomer)
// 	}
// }

// func MustBeRole() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		MustAuthenticateRole(c)
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
// 		if auth.Role != role {
// 			errResponse["error"] = fmt.Sprintf("Unauthorize! you must be %s to access", role)
// 			ctx.JSON(401, errResponse)
// 		} else {
// 			ctx.Next()
// 		}
// 	}
// 	ctx.Abort()
// }

// func MustAuthenticateRole(ctx *gin.Context, role uAuth.Role) {
// 	var errResponse = map[string]interface{}{
// 		"status": "error",
// 	}
// 	if !role.IsValid() {
// 		errResponse["error"] = fmt.Sprintf("Unauthorize! you must be %s to role", role)
// 		ctx.JSON(401, errResponse)
// 	} else {
// 		ctx.Next()
// 	}
// 	ctx.Abort()
// }
