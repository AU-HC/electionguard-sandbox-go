package crypto

import (
	"crypto/rand"
	"electionguard-sandbox-go/constants"
	"math/big"
)

// taken from: https://github.com/AU-HC/elliptic-curve-benchmark-go/blob/master/random/random.go

func GenerateRandomModQ() *big.Int {
	q := constants.GetQ()
	number, err := rand.Int(rand.Reader, q)

	if err != nil {
		panic(err)
	}

	return number
}
