package generation

import (
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
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
		c := *challenges[j]

		var a, b big.Int
		a.Exp(publicKey.G, u, publicKey.P)
		if m == j {
			b.Exp(publicKey.K, u, publicKey.P)
		} else {
			var cMul big.Int
			t := new(big.Int)

			t = t.Add(u, cMul.Mul(big.NewInt(int64(selectionLimit-j)), &c))
			t.Mod(t, publicKey.Q)
			b.Exp(publicKey.K, t, publicKey.P)
		}

		// Calculating the response
		var v big.Int
		v.Sub(u, epsilon.Mul(&epsilon, &c))
		v.Mod(&v, publicKey.Q)

		// Filling the values in we have calculated (note that the c for m == j has to be replaced later)
		cpProofs[j] = models.ChaumPedersenProof{
			Challenge:     c,
			ProofPad:      a,
			ProofData:     b,
			ProofResponse: v,
		}
	}

	// Calculating "true" claim proof
	c := crypto.HMAC(*publicKey.K, 0x21, publicKey.K, alpha, beta)
	cl := new(big.Int)
	cl = cl.Set(c)

	for j, cpProof := range cpProofs {
		if m != j {
			cl.Sub(cl, &cpProof.Challenge)
			cl.Mod(cl, publicKey.Q)
		}
	}

	var v big.Int
	v.Sub(commitments[m], epsilon.Mul(&epsilon, c))
	v.Mod(&v, publicKey.Q)
	cpProofs[m] = models.ChaumPedersenProof{
		Challenge:     *c,
		ProofPad:      cpProofs[m].ProofPad,
		ProofData:     cpProofs[m].ProofData,
		ProofResponse: v,
	}

	// Saving all the proofs into range proof struct
	rangeProof := models.RangeProof{
		Challenge:  *c,
		Proofs:     cpProofs,
		RangeLimit: selectionLimit,
	}

	return rangeProof
}
