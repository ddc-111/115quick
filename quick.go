package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"115Quick_server/internal/config"
	"115Quick_server/internal/handler"
	"115Quick_server/internal/service/auth"
	"115Quick_server/internal/service/cloudtask"
	"115Quick_server/internal/service/download"
	"115Quick_server/internal/service/file"
	"115Quick_server/internal/service/video"
	"115Quick_server/internal/store"
	"115Quick_server/internal/version"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("f", "etc/quick.yaml", "config file path")
	showVersion := flag.Bool("version", false, "show version info")
	flag.Parse()

	if *showVersion {
		info := version.GetInfo()
		fmt.Printf("115Quick Server %s\n", info.Version)
		fmt.Printf("Git Commit: %s\n", info.GitCommit)
		fmt.Printf("Build Time: %s\n", info.BuildTime)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	tokenMgr := auth.NewManager(db)
	apiClient := auth.NewClient(tokenMgr)

	if cfg.Auth115.RefreshToken != "" {
		if err := tokenMgr.SetToken(cfg.Auth115.AccessToken, cfg.Auth115.RefreshToken, 7200); err != nil {
			log.Printf("Warning: Failed to set initial token: %v", err)
		} else {
			if err := tokenMgr.ForceRefresh(); err != nil {
				log.Printf("Warning: Failed to refresh token: %v", err)
			}
		}
	}

	fileSvc := file.NewService(apiClient)
	videoSvc := video.NewService(apiClient, tokenMgr)
	cloudSvc := cloudtask.NewService(apiClient)
	downloadEng := download.NewEngine(tokenMgr, cfg.Auth115.DownloadPath)

	h := handler.NewHandler(tokenMgr, db, fileSvc, videoSvc, cloudSvc, downloadEng)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
