package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hashpw returns the bcrypt hash value of the password
func Password(pw string) (string,error) {
	hashedPw,err := bcrypt.GenerateFromPassword([]byte(pw),bcrypt.DefaultCost)
	if err != nil {
		return "",fmt.Errorf("failed to create Hash password: %w",err)
	}
	return string(hashedPw),nil
}

//CheckPassword - Checks if the provided password is correct or not
func CheckPassword(pw string,hashPw string)error{
	return bcrypt.CompareHashAndPassword([]byte(hashPw),[]byte(pw))
} 