package crypto

import (
	"electionguard-sandbox-go/constants"
	"math/big"
)

// taken from: https://github.com/AU-HC/elliptic-curve-benchmark-go/blob/master/elgamal/elgamal.go

type PublicKey struct {
	G, P, K *big.Int
}

type SecretKey struct {
	G, P, S *big.Int
}

func GenerateKeyPair() (PublicKey, SecretKey) {
	// Fixed G and P to ElectionGuard standard
	var publicKey PublicKey
	var secretKey SecretKey

	g := constants.GetG()
	p := constants.GetP()
	s := GenerateRandomModQ()

	secretKey.G = g
	secretKey.P = p
	secretKey.S = s

	publicKey.G = g
	publicKey.P = p
	publicKey.K = big.NewInt(0).Exp(publicKey.G, s, publicKey.P)

	return publicKey, secretKey
}

func Encrypt(publicKey PublicKey, m, epsilon *big.Int) (*big.Int, *big.Int) {
	var alpha big.Int
	alpha.Exp(publicKey.G, epsilon, publicKey.P) // g^\epsilon mod p

	var kM big.Int
	var kR big.Int
	kM.Exp(publicKey.K, m, publicKey.P)       // K^m mod p
	kR.Exp(publicKey.K, epsilon, publicKey.P) // K^\epsilon mod p

	var beta big.Int
	beta.Mul(&kM, &kR)           // K^m K^\epsilon
	beta.Mod(&beta, publicKey.P) // K^m \cdot K^\epsilon mod p

	return &alpha, &beta // (g^\epsilon, K^m \cdot K^\epsilon mod p)
}
