package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var SessionID string

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("sessionID")
		if err != nil || sessionID == "" {
			ctx.HTML(http.StatusUnauthorized, "admin_login.html", nil)
			ctx.Abort()
			return
		}

		if SessionID != sessionID {
			ctx.HTML(http.StatusUnauthorized, "admin_login.html", nil)
			ctx.Abort()
			return
		}

		ctx.SetCookie(
			"sessionID",
			sessionID,
			3600,
			"/admin",
			"",
			false,
			true,
		)
		ctx.Next()
	}
}
