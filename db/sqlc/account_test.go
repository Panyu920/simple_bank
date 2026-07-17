package db

import (
	"context"
	"log"
	"simple_bank/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *Account {
	params := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	if err != nil {
		log.Fatal("create account err ", err)
	}

	require.Equal(t, account.Owner, params.Owner)
	require.Equal(t, account.Balance, params.Balance)
	require.Equal(t, account.Currency, params.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return &account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
