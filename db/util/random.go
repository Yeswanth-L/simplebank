package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qewuyreposghlgkaznbvbxmcmnvmc"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomInt generates a random number between min and max
func RandomInt(min,max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//Generate a random string og length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i:=0;i< n;i++ {
		char := alphabet[rand.Intn(k)]
		sb.WriteByte(char)
	}
	return sb.String()
}

//Generate random owner name
func RandomOwnerName() string{
   return RandomString(6)
}

//Generates random money
func RandomMoney() int64{
	return RandomInt(100,500)
}

//random currency
func RandomCurreny() string{
	cur := []string{"EUR","USD","CAD"}
	return cur[rand.Intn(len(cur))]
}