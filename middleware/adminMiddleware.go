package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]string
}

var Sessions SessionStore

func InitMiddleware() {
	Sessions = SessionStore{
		sessions: make(map[string]string),
	}
}

func AddSession(sessionID string, username string) {
	Sessions.mu.Lock()
	defer Sessions.mu.Unlock()

	Sessions.sessions[sessionID] = username
}

func GetSession(sessionID string) (string, bool) {
	Sessions.mu.RLock()
	defer Sessions.mu.RUnlock()

	username, exists := Sessions.sessions[sessionID]
	return username, exists
}

func DeleteSession(sessionID string) {
	Sessions.mu.Lock()
	defer Sessions.mu.Unlock()

	delete(Sessions.sessions, sessionID)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("sessionID")
		if err != nil || sessionID == "" {
			ctx.HTML(http.StatusUnauthorized, "admin_login.html", nil)
			ctx.Abort()
			return
		}

		username, exists := GetSession(sessionID)
		if !exists {
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
			true, // secure=true в production с HTTPS
			true,
		)

		ctx.Set("adminLogin", username)

		ctx.Next()
	}
}
