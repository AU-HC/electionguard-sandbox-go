package modular_arithmetic

import (
	"electionguard-sandbox-go/constants"
	"math/big"
)

// taken from: https://github.com/AU-HC/electionguard-verifier-go/blob/master/core/arithmatic.go

func MulP(a, b *big.Int) *big.Int {
	var result big.Int
	p := constants.GetP()

	modOfA := a.Mod(a, p)
	modOfB := b.Mod(b, p)

	// Multiply the two numbers mod p
	result.Mul(modOfA, modOfB)
	result.Mod(&result, p)

	return &result
}

func MulQ(a, b *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	modOfA := a.Mod(a, q)
	modOfB := b.Mod(b, q)

	// Multiply the two numbers mod q
	result.Mul(modOfA, modOfB)
	result.Mod(&result, q)

	return &result
}

func PowP(b, e *big.Int) *big.Int {
	var result big.Int
	p := constants.GetP()

	result.Exp(b, e, p)

	return &result
}

func AddQ(a, b *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	result.Add(b, a)
	result.Mod(&result, q)

	return &result
}

func SubQ(a, b *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	result.Sub(a, b)
	result.Mod(&result, q)

	return &result
}

func ModQ(a *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	result.Mod(a, q)
	return &result
}

func IsValidResidue(a big.Int) bool {
	// Checking the value is in range
	p := constants.GetP()
	q := constants.GetQ()
	zero := big.NewInt(0)
	one := big.NewInt(1)

	valueIsAboveOrEqualToZero := zero.Cmp(&a) <= 0
	valueIsSmallerThanP := p.Cmp(&a) == 1
	valueIsInRange := valueIsSmallerThanP && valueIsAboveOrEqualToZero // a is in [0, P)

	validResidue := PowP(&a, q).Cmp(one) == 0 // a^q mod p == 1

	return valueIsInRange && validResidue
}

func IsInRange(a big.Int) bool {
	q := constants.GetQ()
	zero := big.NewInt(0)

	valueIsAboveOrEqualToZero := zero.Cmp(&a) <= 0
	valueIsSmallerThanP := q.Cmp(&a) > 0

	return valueIsAboveOrEqualToZero && valueIsSmallerThanP
}
