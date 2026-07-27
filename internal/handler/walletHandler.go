package handler

import (
	"errors"
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

type TransactionReq struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Ref    string  `json:"reference" binding:"required"`
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

func (h *walletHandler) CreditBalance(c *gin.Context) {
	var req TransactionReq

	val, exists := c.Get("user")
	if !exists {
		h.log.Warn("user not found context")
		response.JSON(c, http.StatusUnauthorized, "user not found context", nil)
		c.Abort()
		return
	}

	tokenData := val.(*domain.TokenData)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed parse request : ", zap.Error(err))
		response.JSON(c, http.StatusBadRequest, "failed parse request", nil)
		return
	}

	_, err := h.uc.Credit(c.Request.Context(), tokenData.UserID, req.Amount, req.Ref)

	if err != nil {
		if errors.Is(err, domain.ErrDuplicateReference) {
			h.log.Warn("transaction ref duplicate", zap.Error(err))
			response.JSON(c, http.StatusConflict, "transaction ref duplicated", nil)
			c.Abort()
			return
		}

		if errors.Is(err, domain.ErrInsufficientBalance) {
			h.log.Warn("insufficient balance", zap.Error(err))
			response.JSON(c, http.StatusBadRequest, "insufficient balance", nil)
			c.Abort()
			return
		}

		if errors.Is(err, domain.ErrRecordNotFound) {
			h.log.Warn("wallet not found", zap.Error(err))
			response.JSON(c, http.StatusNotFound, "wallet not found", nil)
			c.Abort()
			return
		}

		h.log.Error("internal server error", zap.Error(err))
		response.JSON(c, http.StatusInternalServerError, "internal server error", nil)
		c.Abort()
		return
	}

	response.JSON(c, http.StatusOK, "success credit balance", nil)

}

func (h *walletHandler) DebitBalance(c *gin.Context) {
	var req TransactionReq

	val, exists := c.Get("user")
	if !exists {
		h.log.Warn("user not found context")
		response.JSON(c, http.StatusUnauthorized, "user not found context", nil)
		c.Abort()
		return
	}

	tokenData := val.(*domain.TokenData)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed parse request : ", zap.Error(err))
		response.JSON(c, http.StatusBadRequest, "failed parse request", nil)
		return
	}

	_, err := h.uc.Debit(c.Request.Context(), tokenData.UserID, req.Amount, req.Ref)

	if err != nil {
		if errors.Is(err, domain.ErrDuplicateReference) {
			h.log.Warn("transaction ref duplicate", zap.Error(err))
			response.JSON(c, http.StatusConflict, "transaction ref duplicated", nil)
			c.Abort()
			return
		}

		if errors.Is(err, domain.ErrInsufficientBalance) {
			h.log.Warn("insufficient balance", zap.Error(err))
			response.JSON(c, http.StatusBadRequest, "insufficient balance", nil)
			c.Abort()
			return
		}

		if errors.Is(err, domain.ErrRecordNotFound) {
			h.log.Warn("wallet not found", zap.Error(err))
			response.JSON(c, http.StatusNotFound, "wallet not found", nil)
			c.Abort()
			return
		}

		h.log.Error("internal server error", zap.Error(err))
		response.JSON(c, http.StatusInternalServerError, "internal server error", nil)
		c.Abort()
		return
	}

	response.JSON(c, http.StatusOK, "success debit balance", nil)

}

func (h *walletHandler) GetBalance(c *gin.Context) {
	val, exists := c.Get("user")
	if !exists {
		h.log.Warn("user not found context")
		response.JSON(c, http.StatusUnauthorized, "user unauthorized", nil)
		return
	}

	tokenData := val.(*domain.TokenData)

	wallet, err := h.uc.GetBalance(c.Request.Context(), tokenData.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			h.log.Warn("wallet not found", zap.Error(err))
			response.JSON(c, http.StatusNotFound, "wallet not found", nil)
			return
		}

		h.log.Error("internal server error", zap.Error(err))
		response.JSON(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.JSON(c, http.StatusOK, "success get balance", wallet)
}

func (h *walletHandler) GetTransactionHistory(c *gin.Context) {
	val, exists := c.Get("user")

	if !exists {
		h.log.Warn("user not found context")
		response.JSON(c, http.StatusUnauthorized, "user unauthorized", nil)
		return
	}

	tokenData := val.(*domain.TokenData)

	var param domain.WalletHistoryParam

	if err := c.ShouldBindQuery(&param); err != nil {
		h.log.Error("failed parse params", zap.Error(err))
		response.JSON(c, http.StatusBadRequest, "failed parse params", nil)
		return
	}

	if param.Page <= 0 {
		param.Page = 1
	}

	if param.Limit <= 0 {
		param.Limit = 10
	}

	if param.Type != "" && param.Type != domain.WalletTransactionTypeCredit && param.Type != domain.WalletTransactionTypeDebit {
		h.log.Warn("invalid wallet transaction type")
		response.JSON(c, http.StatusBadRequest, "invalid wallet transaction type", nil)
		return
	}

	listTrx, err := h.uc.GetHistory(c.Request.Context(), tokenData.UserID, param)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			h.log.Warn("wallet not found", zap.Error(err))
			response.JSON(c, http.StatusNotFound, "wallet not found", nil)
			return
		}

		h.log.Error("internal server error", zap.Error(err))
		response.JSON(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.JSON(c, http.StatusOK, "success get transactions", listTrx)

}
