package router

import (
	"github.com/gin-gonic/gin"
	aihandler "github.com/leoferamos/aroma-sense/internal/handler/ai"
	chathandler "github.com/leoferamos/aroma-sense/internal/handler/chat"
)

// AIRoutes registers AI chat related routes.
func AIRoutes(r *gin.Engine, aiHandler *aihandler.AIHandler, chatHandler *chathandler.ChatHandler) {
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
