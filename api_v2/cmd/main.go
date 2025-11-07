package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/server"
	"github.com/joho/godotenv"
)

const banner = `
 ██████╗██╗███╗   ██╗███████╗    ██╗   ██╗██████╗ 
██╔════╝██║████╗  ██║██╔════╝    ██║   ██║╚════██╗
██║     ██║██╔██╗ ██║█████╗      ██║   ██║ █████╔╝
██║     ██║██║╚██╗██║██╔══╝      ╚██╗ ██╔╝██╔═══╝ 
╚██████╗██║██║ ╚████║███████╗     ╚████╔╝ ███████╗
 ╚═════╝╚═╝╚═╝  ╚═══╝╚══════╝      ╚═══╝  ╚══════╝
                                                  
 🎬 CineVerse API v2.0.0 - Authentication Sprint 1
`

func main() {
	fmt.Print(banner)

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup logger
	logger := middleware.SetupLogger(cfg.Server.Environment)
	slog.SetDefault(logger)

	// Create server
	srv, err := server.NewServer(cfg, logger)
	if err != nil {
		logger.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("🚀 Server started successfully!",
		"address", cfg.Server.GetServerAddress(),
		"environment", cfg.Server.Environment)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Server is shutting down gracefully...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Stop(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("✅ Server stopped successfully")
}
