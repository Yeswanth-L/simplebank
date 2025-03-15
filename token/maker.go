package token

import "time"

//Maker is an interface for managing token (like JWT, Paesto...)
//Any struct that implements both methods is considered to satisfy the Maker interface.
type Maker interface {
	//CreateToken - creates a token for the specified username & duration
	CreateToken(username string, duration time.Duration) (string,error)

	//VerifyToken - Checks if the tokem is valid or not
	VerifyToken(token string) (*Payload,error)
}