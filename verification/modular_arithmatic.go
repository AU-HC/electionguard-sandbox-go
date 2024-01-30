package verification

import (
	"electionguard-sandbox-go/constants"
	"math/big"
)

// taken from: https://github.com/AU-HC/electionguard-verifier-go/blob/master/core/arithmatic.go

func mulP(a, b *big.Int) *big.Int {
	var result big.Int
	p := constants.GetP()

	modOfA := a.Mod(a, p)
	modOfB := b.Mod(b, p)

	// Multiply the two numbers mod p
	result.Mul(modOfA, modOfB)
	result.Mod(&result, p)

	return &result
}

func mulQ(a, b *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	modOfA := a.Mod(a, q)
	modOfB := b.Mod(b, q)

	// Multiply the two numbers mod q
	result.Mul(modOfA, modOfB)
	result.Mod(&result, q)

	return &result
}

func powP(b, e *big.Int) *big.Int {
	var result big.Int
	p := constants.GetP()

	result.Exp(b, e, p)

	return &result
}

func addQ(a, b *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	result.Add(b, a)
	result.Mod(&result, q)

	return &result
}

func subQ(a, b *big.Int) *big.Int {
	var result big.Int
	q := constants.GetQ()

	result.Sub(a, b)
	result.Mod(&result, q)

	return &result
}

func isValidResidue(a big.Int) bool {
	// Checking the value is in range
	p := constants.GetP()
	q := constants.GetQ()
	zero := big.NewInt(0)
	one := big.NewInt(1)

	valueIsAboveOrEqualToZero := zero.Cmp(&a) <= 0
	valueIsSmallerThanP := p.Cmp(&a) == 1
	valueIsInRange := valueIsSmallerThanP && valueIsAboveOrEqualToZero // a is in [0, P)

	validResidue := powP(&a, q).Cmp(one) == 0 // a^q mod p == 1

	return valueIsInRange && validResidue
}

func isInRange(a big.Int) bool {
	q := constants.GetQ()
	zero := big.NewInt(0)

	valueIsAboveOrEqualToZero := zero.Cmp(&a) <= 0
	valueIsSmallerThanP := q.Cmp(&a) > 0

	return valueIsAboveOrEqualToZero && valueIsSmallerThanP
}
