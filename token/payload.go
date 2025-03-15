package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = fmt.Errorf("token is invalid")
	ErrExpiredToken =  fmt.Errorf("token has expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id`
	Username  string    `json:"username"`
	IssuedAt  time.Time    `json:"issued_at"`
	ExpiredAt time.Time   `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload,error){
	tokenID,err := uuid.NewRandom()  //it generates the tokenID for the token
	if err != nil {
		return nil,err
	}

	payload := &Payload{
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(duration)),
	}
	return payload,nil
}

//Valid - checks the token is valid or not
func(payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt){
		return ErrExpiredToken
	}
	return nil
}