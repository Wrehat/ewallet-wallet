package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Wrehat/ewallet-wallet/internal/config"
	"github.com/Wrehat/ewallet-wallet/internal/handler"
	"github.com/Wrehat/ewallet-wallet/internal/middleware"
	"github.com/Wrehat/ewallet-wallet/internal/repository"
	"github.com/Wrehat/ewallet-wallet/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ServeHTTP(ctx context.Context, cfg *config.AppConfig, log *zap.Logger, db *gorm.DB) {

	hCheckRepo := repository.NewHealthRepo(db)
	hCheckUc := usecase.NewHealthUsecase(hCheckRepo)
	hCheckHndlr := handler.NewHealthHandler(hCheckUc)

	walletRepo := repository.NewWalletRepo(db)
	walletUc := usecase.NewWalletUsecase(walletRepo, log)
	walletHndlr := handler.NewWalletHandler(walletUc, log)

	r := gin.New()
	r.Use(gin.Recovery())

	umsGateway, err := repository.NewUMSGateway(cfg.UmsGrpcHost)
	if err != nil {
		log.Fatal("failed to create UMS gateway : ", zap.Error(err))
	}

	authMiddleware := middleware.AuthMiddleware(umsGateway, log)

	apiV1 := r.Group("/api/v1")
	walletGroup := apiV1.Group("/wallets")
	{
		walletGroup.POST("/", walletHndlr.CreateWallet)
		walletAuth := walletGroup.Use(authMiddleware)
		{
			walletAuth.PUT("/credit", walletHndlr.CreditBalance)
			walletAuth.PUT("/debit", walletHndlr.DebitBalance)
		}
	}

	r.GET("/health", hCheckHndlr.HealthCheck)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("failed shutdown server : ", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error HTTP server", zap.Error(err))
	}

}
