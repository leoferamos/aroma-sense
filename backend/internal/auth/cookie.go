package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetRefreshTokenCookie sets the refresh token cookie
func SetRefreshTokenCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // Allow cross-site cookies
	}
	http.SetCookie(c.Writer, cookie)
}

// ClearRefreshTokenCookie removes the refresh token cookie
func ClearRefreshTokenCookie(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(c.Writer, cookie)
}
