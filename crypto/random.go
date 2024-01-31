package crypto

import (
	"crypto/rand"
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/models"
)

// taken from: https://github.com/AU-HC/elliptic-curve-benchmark-go/blob/master/random/random.go

func GenerateRandomModQ() *models.BigInt {
	q := constants.GetQ()
	number, err := rand.Int(rand.Reader, &q.Int)

	if err != nil {
		panic(err)
	}

	return &models.BigInt{Int: *number}
}
