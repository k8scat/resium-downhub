package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8scat/downhub/config"
	"github.com/k8scat/downhub/util"
)

type AuthHeaders struct {
	Token string `header:"token"`
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var headers AuthHeaders
		if err := c.ShouldBindHeader(&headers); err != nil {
			c.Abort()
		}
		if headers.Token != config.Config.Token {
			c.JSON(util.JSONResponse(http.StatusUnauthorized, "unauthorized", nil))
			c.Abort()
		}
		c.Next()
	}
}
