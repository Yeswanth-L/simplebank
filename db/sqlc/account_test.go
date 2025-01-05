package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
)

func createRandomAccount(t *testing.T) Account{
	params := CreateAccountParams{
		Owner:    util.RandomOwnerName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurreny(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t,params.Owner,account.Owner)
	require.Equal(t,params.Balance,account.Balance)
	require.Equal(t,params.Currency,account.Currency)
	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(),acc1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,acc2)
	require.Equal(t, acc1.Balance,acc2.Balance)
	require.WithinDuration(t,acc1.CreatedAt,acc2.CreatedAt,time.Second)
}

func TestUpdateAccount(t *testing.T){
	acc1 := createRandomAccount(t)
	params := UpdateAccountParams{
		ID: acc1.ID,
		Balance: acc1.Balance,
	}
	testQueries.UpdateAccount(context.Background(),params)
}

func TestDeleteAccount(t *testing.T){
	acc1 := createRandomAccount(t)
	testQueries.DeleteAccount(context.Background(),acc1.ID)
}

func TestListAccounts(t *testing.T){
	for i:=0 ; i<10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams {
		Limit: 5,
		Offset: 5,
	}

	testQueries.ListAccounts(context.Background(),args)

}