package handler

import (
	"net/http"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"github.com/Wrehat/ewallet-wallet/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type walletHandler struct {
	uc  domain.WalletUsecase
	log *zap.Logger
}

type CreateWalletReq struct {
	UserID int `json:"user_id" binding:"required"`
}

func NewWalletHandler(uc domain.WalletUsecase, log *zap.Logger) *walletHandler {
	return &walletHandler{
		uc:  uc,
		log: log,
	}
}

func (h *walletHandler) CreateWallet(c *gin.Context) {
	var req CreateWalletReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		response.JSON(c, http.StatusBadRequest, "failed to parse request", nil)
		return
	}

	if req.UserID == 0 {
		response.JSON(c, http.StatusBadRequest, "user_id is required", nil)
		return
	}

	err := h.uc.Create(c.Request.Context(), &domain.Wallet{
		UserID: req.UserID,
	})
	if err != nil {
		h.log.Error("internal server error", zap.Error(err))
		response.JSON(c, http.StatusInternalServerError, "failed create wallet", nil)
		return
	}

	response.JSON(c, http.StatusCreated, "success create wallet", nil)

}
