package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// AIRoutes registers AI chat related routes.
func AIRoutes(r *gin.Engine, aiHandler *handler.AIHandler, chatHandler *handler.ChatHandler) {
	if aiHandler == nil && chatHandler == nil {
		return
	}
	grp := r.Group("/ai")
	{
		if aiHandler != nil {
			grp.POST("/recommend", aiHandler.Recommend)
		}
		if chatHandler != nil {
			grp.POST("/chat", chatHandler.Chat)
		}
	}
}
