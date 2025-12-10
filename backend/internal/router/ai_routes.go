package router

import (
	"github.com/gin-gonic/gin"
	chathandler "github.com/leoferamos/aroma-sense/internal/handler/chat"
)

// AIRoutes registers AI chat related routes.
func AIRoutes(r *gin.Engine, aiHandler *chathandler.AIHandler, chatHandler *chathandler.ChatHandler) {
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
