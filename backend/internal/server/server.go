package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leoferamos/aroma-sense/internal/bootstrap"
	"github.com/leoferamos/aroma-sense/internal/job"
	"github.com/leoferamos/aroma-sense/internal/router"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// StartServer starts background jobs, the HTTP server and manages graceful shutdown.
func StartServer(app *bootstrap.AppComponents, addr string) error {
	// Initialize and start LGPD jobs
	autoConfirmJob := job.NewAutoConfirmJob(
		app.Repos.UserRepo,
		app.Services.LgpdService,
		app.Services.AuditLogService,
	)
	autoConfirmJob.Start()

	cleanupJob := job.NewDataCleanupJob(
		app.Repos.UserRepo,
		app.Services.LgpdService,
		app.Services.AuditLogService,
	)
	cleanupJob.Start()

	// Setup router with all handlers
	r := router.SetupRouter(app.Handlers)

	// Swagger docs route
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Shutdown integrations
	bootstrap.ShutdownIntegrations(ctx)

	log.Println("Server exited cleanly")
	return nil
}
