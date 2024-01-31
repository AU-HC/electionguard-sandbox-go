package verification

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	mod "electionguard-sandbox-go/modular_arithmetic"
	"math/big"
)

func VerifyStep5(ballots []models.Ballot, publicKey crypto.PublicKey) {
	g := publicKey.G
	k := publicKey.K

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				alpha := selection.Ciphertext.Pad
				beta := selection.Ciphertext.Data

				toBeHashed := []interface{}{*publicKey.K, selection.Ciphertext.Pad, selection.Ciphertext.Data}
				computedC := big.NewInt(0)
				for j, proof := range selection.Proof.Proofs {
					cj := proof.Challenge
					computedC = mod.AddQ(computedC, &cj)

					vj := proof.ProofResponse
					wj := mod.SubQ(&vj, mod.MulQ(big.NewInt(int64(j)), &cj))

					aj := mod.MulP(mod.PowP(g, &vj), mod.PowP(&alpha, &cj))
					bj := mod.MulP(mod.PowP(k, wj), mod.PowP(&beta, &cj))

					toBeHashed = append(toBeHashed, *aj)
					toBeHashed = append(toBeHashed, *bj)

					verify("(5.B) ...", mod.IsInRange(cj))
					verify("(5.C) ...", mod.IsInRange(vj))
				}

				c := crypto.HMAC(constants.GetExtendedBaseHash(), 0x21, toBeHashed...)

				verify("(5.A) alpha is in Z_p^r", mod.IsValidResidue(alpha))
				verify("(5.A) beta is in Z_p^r", mod.IsValidResidue(beta))
				verify("(5.D) challenge is computed correctly", c.Cmp(computedC) == 0)
			}
		}
	}
}
