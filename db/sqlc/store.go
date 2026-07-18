package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, params TransferParams) (TransferResult, error)
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 开启事务

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v , rb err %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SqlStore) TransferTx(ctx context.Context, params TransferParams) (TransferResult, error) {
	var result TransferResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// 创建交易题目
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})
		if err != nil {
			return err
		}

		// 创建entryentry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})

		if err != nil {
			return err
		}

		// 创建收款方条目

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})

		if err != nil {
			return err
		}

		// 更新双方账号余额
		if params.FromAccountID < params.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q,
				params.FromAccountID, -params.Amount, params.ToAccountID, params.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q,
				params.ToAccountID, params.Amount, params.FromAccountID, -params.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, ID1 int64, amount1 int64, ID2 int64, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     ID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     ID2,
		Amount: amount2,
	})

	if err != nil {
		return
	}
	return
}
