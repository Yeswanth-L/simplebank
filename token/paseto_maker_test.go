package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomString(6)
	duration := time.Minute

	issueddAt := time.Now()
	expiredAt := issueddAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueddAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoExpiredToken(t *testing.T){
	maker,err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t,err)

	username := util.RandomString(6)
	duration := time.Minute

	token,err := maker.CreateToken(username,-duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)

	payload,err := maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,ErrExpiredToken.Error())
	require.Nil(t,payload)
}