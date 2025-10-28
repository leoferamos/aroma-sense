package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(c *gin.Context, token string) {
	const defaultExpiryMins = 15
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   defaultExpiryMins * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // Allow cross-site cookies
	}
	http.SetCookie(c.Writer, cookie)
}

// ClearAuthCookie removes the authentication cookie
func ClearAuthCookie(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(c.Writer, cookie)
}
