package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(c *gin.Context, token string) {
	domain := os.Getenv("COOKIE_DOMAIN")

	const defaultExpiryMins = 15
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Domain:   domain,
		MaxAge:   defaultExpiryMins * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)
}

// ClearAuthCookie removes the authentication cookie
func ClearAuthCookie(c *gin.Context) {
	domain := os.Getenv("COOKIE_DOMAIN")

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1, // Delete cookie
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)
}
