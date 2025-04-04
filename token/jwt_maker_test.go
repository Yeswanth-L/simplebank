package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
)
func TestJWTMaker(t *testing.T){
	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)

	username := util.RandomString(6)
	duration := time.Minute

	issueddAt := time.Now()
	expiredAt  := issueddAt.Add(duration)

	token,payload,err := maker.CreateToken(username,duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	payload,err = maker.VerifyToken(token)
	require.NoError(t,err)
	require.NotEmpty(t,payload)
	require.NotZero(t,payload.ID)
	require.Equal(t,username,payload.Username)
	require.WithinDuration(t,issueddAt,payload.IssuedAt,time.Second)
	require.WithinDuration(t,expiredAt,payload.ExpiredAt,time.Second)
}

func TestExpiredToken(t *testing.T){
	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)

	username := util.RandomString(6)
	duration := time.Minute

	token,payload,err := maker.CreateToken(username,-duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	payload,err = maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,ErrExpiredToken.Error())
	require.Nil(t,payload)
}

func TestJWTTokenAlgNone(t *testing.T){
	payload,err := NewPayload(util.RandomOwnerName(),time.Minute)
	require.NoError(t,err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone,payload)
	token,err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t,err)
	
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)

	payload,err = maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,ErrInvalidToken.Error())
	require.Nil(t,payload)
}