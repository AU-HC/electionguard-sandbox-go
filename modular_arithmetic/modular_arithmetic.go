package modular_arithmetic

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/models"
	"math/big"
)

// taken from: https://github.com/AU-HC/electionguard-verifier-go/blob/master/core/arithmatic.go

func MulP(a, b *models.BigInt) *models.BigInt {
	var result models.BigInt
	p := constants.GetP()

	modOfA := a.Mod(&a.Int, &p.Int)
	modOfB := b.Mod(&b.Int, &p.Int)

	// Multiply the two numbers mod p
	result.Mul(modOfA, modOfB)
	result.Mod(&result.Int, &p.Int)

	return &result
}

func MulQ(a, b *models.BigInt) *models.BigInt {
	var result models.BigInt
	q := constants.GetQ()

	modOfA := a.Mod(&a.Int, &q.Int)
	modOfB := b.Mod(&b.Int, &q.Int)

	// Multiply the two numbers mod q
	result.Mul(modOfA, modOfB)
	result.Mod(&result.Int, &q.Int)

	return &result
}

func PowP(b, e *models.BigInt) *models.BigInt {
	var result models.BigInt
	p := constants.GetP()

	result.Exp(&b.Int, &e.Int, &p.Int)

	return &result
}

func AddQ(a, b *models.BigInt) *models.BigInt {
	var result models.BigInt
	q := constants.GetQ()

	result.Add(&b.Int, &a.Int)
	result.Mod(&result.Int, &q.Int)

	return &result
}

func SubQ(a, b *models.BigInt) *models.BigInt {
	var result models.BigInt
	q := constants.GetQ()

	result.Sub(&a.Int, &b.Int)
	result.Mod(&result.Int, &q.Int)

	return &result
}

func ModQ(a *models.BigInt) *models.BigInt {
	var result models.BigInt
	q := constants.GetQ()

	result.Mod(&a.Int, &q.Int)
	return &result
}

func IsValidResidue(a models.BigInt) bool {
	// Checking the value is in range
	p := constants.GetP()
	q := constants.GetQ()
	zero := big.NewInt(0)
	one := big.NewInt(1)

	valueIsAboveOrEqualToZero := zero.Cmp(&a.Int) <= 0
	valueIsSmallerThanP := p.Cmp(&a.Int) == 1
	valueIsInRange := valueIsSmallerThanP && valueIsAboveOrEqualToZero // a is in [0, P)

	validResidue := PowP(&a, models.MakeBigIntFromByteArray(q.Bytes())).Cmp(one) == 0 // a^q mod p == 1

	return valueIsInRange && validResidue
}

func IsInRange(a big.Int) bool {
	q := constants.GetQ()
	zero := big.NewInt(0)

	valueIsAboveOrEqualToZero := zero.Cmp(&a) <= 0
	valueIsSmallerThanP := q.Cmp(&a) > 0

	return valueIsAboveOrEqualToZero && valueIsSmallerThanP
}
