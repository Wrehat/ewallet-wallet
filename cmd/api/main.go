package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Wrehat/ewallet-wallet/internal/config"
	"github.com/Wrehat/ewallet-wallet/internal/database"
	"github.com/Wrehat/ewallet-wallet/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Setup config untuk data environment
	cfg := config.SetupConfig()

	// Setup logger untuk tracing
	zapLog, err := logger.SetupLogger(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed init logger : %v", err)
	}
	defer zapLog.Sync()

	// Setup database
	db, err := database.SetupDB(cfg.DSN(), zapLog)
	if err != nil {
		zapLog.Fatal("failed setup database", zap.Error(err))
	}
	sqlDb, err := db.DB()
	if err == nil {
		defer sqlDb.Close()
	}

	// Get signal for shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run server with waitgroup n goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ServeGRPC(ctx, cfg, zapLog)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ServeHTTP(ctx, cfg, zapLog, db)
	}()

	zapLog.Info("server running...")

	// Wait for shutdown signal
	<-ctx.Done()
	wg.Wait()
	zapLog.Info("server gracful shutdown.")

}
