package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T){
	pw := RandomString(6)

	hashPw1, err := Password(pw)
	require.NoError(t,err)
	require.NotEmpty(t,hashPw1)

	hashPw2, err := Password(pw) //It wont provide the same hashvalue for the same pw. Everytime it will produce diff. hashvalue.
	require.NoError(t,err)
	require.NotEmpty(t,hashPw2)
	require.NotEqual(t,hashPw1,hashPw2) 

	err1 := CheckPassword(pw,hashPw1)
	require.NoError(t,err1)

	pw1 := RandomString(6)
	err2 := CheckPassword(pw1,hashPw1)
	require.EqualError(t,err2,bcrypt.ErrMismatchedHashAndPassword.Error())

}