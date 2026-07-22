package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, ok := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !ok {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload.Username != fromAccount.Owner {
		err := errors.New("account does not belong to authcatied user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if req.Amount > fromAccount.Balance {
		err := errors.New("余额不足")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, ok := server.validAccount(ctx, req.ToAccountID, req.Currency); !ok {
		return
	}

	arg := db.TransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (*db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return nil, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return nil, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("accountID [%d] currency error : %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return nil, false
	}

	return &account, true
}
