package db

import (
	"context"
	"simple_bank/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	params := CreateUserParams{
		Username:       utils.RandomOwner(),
		FullName:       utils.RandomOwner(),
		HashedPassword: "123456",
		Email:          utils.RandomEmail(),
	}

	account, err := testQueries.CreateUser(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Username, params.Username)
	require.Equal(t, account.FullName, params.FullName)
	require.Equal(t, account.HashedPassword, params.HashedPassword)
	require.NotZero(t, account.CreatedAt)

	return &account
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
