package middleware

import (
	"GamesAPI/src/services"
	"github.com/gin-gonic/gin"
)

///Code inspired from article:
// https://sosedoff.com/2014/12/21/gin-middleware.html?fbclid=IwAR2n3Kvdb1kgYitha4ur2q1bzrN5yukIESg2BKHEe6AoZdvanaZuYqpjwmE
// (27 oct 2020)

func InitApiToken(r *gin.Engine) {
	r.Use(MiddlewareHandler)
}

// Message Contexte Error
func ErrorMessageTypeCode(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"Error": message})
}

//Token  Authentificator  handler
func MiddlewareHandler(c *gin.Context) {
	/// token := c.Request.FormValue("api_token") //.header

	token := c.Request.Header.Get("x-api-key")

	if token == "" {
		ErrorMessageTypeCode(c, 400, "API token required")
		return
	}
	resultValidate, err := services.TokenService.ValidateToken(token)
	if err != nil {
		ErrorMessageTypeCode(c, 500, err.Error())
		return
	}

	if !resultValidate {
		ErrorMessageTypeCode(c, 401, "Invalid API token")
		return
	}

	c.Next()
}
