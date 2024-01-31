package generation

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	mod "electionguard-sandbox-go/modular_arithmetic"
	"math/big"
)

func generateRangeProofFromEncryptionAndNonce(alpha, beta, epsilon big.Int, publicKey crypto.PublicKey, selectionLimit, m int) models.RangeProof {
	// Generate random challenges
	cpProofs := make([]models.ChaumPedersenProof, selectionLimit+1) // need R+1 proofs
	challenges := getRandomNumbersModQ(selectionLimit + 1)          // need R+1 challenges
	commitments := getRandomNumbersModQ(selectionLimit + 1)         // need R+1 commitments

	// Generate range proof of this encryption
	for j := 0; j <= selectionLimit; j++ {
		u := commitments[j]
		cj := *challenges[j]

		a := mod.PowP(publicKey.G, u)
		b := calculateBCommitment(m, j, publicKey, *u, cj)

		// Calculating the response
		v := mod.SubQ(u, mod.MulQ(&epsilon, &cj))

		// Filling the values in we have calculated (note that the c for m == j has to be replaced later)
		cpProofs[j] = models.ChaumPedersenProof{
			Challenge:     cj,
			ProofPad:      *a,
			ProofData:     *b,
			ProofResponse: *v,
		}
	}

	// Calculating "true" claim proof
	xd := []interface{}{publicKey.K, alpha, beta}
	for i := 0; i < len(cpProofs); i++ {
		xd = append(xd, cpProofs[i].ProofPad)
		xd = append(xd, cpProofs[i].ProofData)
	}

	c := crypto.HMAC(constants.GetExtendedBaseHash(), 0x21, xd...)
	cl := new(big.Int)
	cl = cl.Set(c)

	for j, cpProof := range cpProofs {
		if m != j {
			cl = mod.SubQ(cl, &cpProof.Challenge)
		}
	}

	var v *big.Int
	v = mod.SubQ(commitments[m], mod.MulQ(&epsilon, cl))

	cpProofs[m] = models.ChaumPedersenProof{
		Challenge:     *cl,
		ProofPad:      cpProofs[m].ProofPad,
		ProofData:     cpProofs[m].ProofData,
		ProofResponse: *v,
	}

	// Saving all the proofs into range proof struct
	rangeProof := models.RangeProof{
		Challenge:  *c,
		Proofs:     cpProofs,
		RangeLimit: selectionLimit,
	}

	return rangeProof
}

func calculateBCommitment(m, j int, publicKey crypto.PublicKey, u, c big.Int) *big.Int {
	if m == j {
		return mod.PowP(publicKey.K, &u)
	} else {
		t := mod.AddQ(&u, mod.MulQ(big.NewInt(int64(m-j)), &c))
		return mod.PowP(publicKey.K, t)
	}
}
