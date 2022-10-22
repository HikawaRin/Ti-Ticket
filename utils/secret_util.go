package utils

import (
	"math/rand"
	"time"
)

var (
	min int = '0'
	max int = 'Z'
)

func GetSecret(len int) string {
	rand.Seed(time.Now().Unix())
	secret := ""
	for i := 0; i < len; i++ {
		secret = secret + string(byte(rand.Intn(max-min)+min))
	}
	return secret
}
